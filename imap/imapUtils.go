package imap

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/pierelucas/atlantr-extreme/utils"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	strip "github.com/grokify/html-strip-tags-go"
	"golang.org/x/net/proxy"
)

// Imaper --
type Imaper struct {
	sendersToSave  []string
	useSocks       bool
	proccessEmails bool
}

// NewImap --
func NewImap(sendersToSave []string, useSocks, proccessEmails bool) *Imaper {
	return &Imaper{
		sendersToSave:  sendersToSave,
		useSocks:       useSocks,
		proccessEmails: proccessEmails,
	}
}

// IMAPutil --
func (im *Imaper) IMAPutil(socksAddr string, addr string, emailUser string, emailPassword string) (bool, error) {
	var c *client.Client

	if im.useSocks {
		cc, err := connextWithSocks5(socksAddr, addr)
		if err != nil {
			return false, err
		}
		c = cc
	} else {
		cc, err := connectTLS(addr)
		if err != nil {
			return false, err
		}
		c = cc
	}

	defer c.Logout()

	err := c.Login(emailUser, emailPassword)
	if err != nil {
		return false, err
	}

	if im.proccessEmails {
		inboxProcessing(c, im.sendersToSave)
	}
	return true, nil
}

func inboxProcessing(c *client.Client, sendersToSave []string) {
	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 100)
	done := make(chan error, 1)

	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	var mailboxList []string
	for m := range mailboxes {
		mailboxList = append(mailboxList, m.Name)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	findMatchingSenders(mailboxList, c, sendersToSave)
}

func findMatchingSenders(mailboxList []string, c *client.Client, sendersToSave []string) {
	for i := range mailboxList {
		mbox, err := c.Select(mailboxList[i], true)
		if err != nil {
			//log.Println(err)
			continue
		}

		if mbox.Messages < 1 {
			continue
		}

		from := uint32(1)
		to := mbox.Messages
		const maxMessagesToGet = 100
		if mbox.Messages > maxMessagesToGet {
			to = from + maxMessagesToGet - 1
		}
		seqset := new(imap.SeqSet)
		seqset.AddRange(from, to)

		// Get the whole message body
		var section imap.BodySectionName
		items := []imap.FetchItem{section.FetchItem()}

		messages := make(chan *imap.Message, maxMessagesToGet+1)
		done := make(chan error, 1)

		go func() {
			done <- c.Fetch(seqset, items, messages)
		}()

		if err := <-done; err != nil {
			log.Println("fetching error:", err)
		}

		for msg := range messages {
			if msg == nil {
				continue
			}

			r := msg.GetBody(&section)
			if r == nil {
				continue
			}

			// Create a new mail reader
			mr, err := mail.CreateReader(r)
			if err != nil {
				//log.Println("createReader", err)
				continue
			}

			// Print some info about the message
			header := mr.Header

			from, err := header.AddressList("From")
			if err != nil {
				//log.Println("header", err)
				continue
			}

			if len(from) < 1 {
				continue
			}

			for i := range sendersToSave {
				if strings.Contains(from[0].Address, sendersToSave[i]) {
					//log.Println(from[0].Address)
					// Process each message's part
					for {
						p, err := mr.NextPart()
						if err == io.EOF {
							break
						} else if err != nil {
							//log.Println("part", err)
							continue
						}

						switch h := p.Header.(type) {
						case *mail.InlineHeader:
							// This is the message's text (can be plain-text or HTML)
							b, err := ioutil.ReadAll(p.Body)
							if err != nil {
								log.Println("read body", err)
							} else {
								bb := strip.StripTags(string(b))
								//singleSpacePattern := regexp.MustCompile(`\s+`)
								//bbb := singleSpacePattern.ReplaceAllString(bb, " ")
								//log.Println(bbb)

								appendStringToFile(sendersToSave[i], bb)
							}
						case *mail.AttachmentHeader:
							// This is an attachment
							_, err := h.Filename()
							if err != nil {
								log.Println(err)
							}
							//log.Println("Got attachment: %v", filename)
						}
					}
				}
			}
		}
	}
}

func connectTLS(addr string) (*client.Client, error) {
	// Connect to server
	c, err := client.DialTLS(addr, nil)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func connextWithSocks5(socksAddr string, addr string) (*client.Client, error) {
	// create a socks5 dialer
	dialer, err := proxy.SOCKS5("tcp", socksAddr, nil, proxy.Direct)
	if err != nil {
		return nil, err
	}

	// Connect to server
	c, err := client.DialWithDialerTLS(dialer, addr, nil)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func appendStringToFile(fileName string, text string) {
	fileName += ".txt" // adding .txt to file cause maybe some noobish windows users have problems to open the file

	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	utils.CheckError(err)
	defer f.Close()

	divider := "\r\n\r\n$#------------this separates emails---------------#$\r\n\r\n"

	_, err = f.WriteString(text + divider)
	utils.CheckError(err)
}
