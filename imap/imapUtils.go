package imap

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
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
	SendersToSave    []string
	UseSocks         bool
	ProccessEmails   bool
	DownloadEmails   bool
	MaxMessagesToGet uint32
	MatcherBaseDir   string
}

// NewImap --
func NewImap(sendersToSave []string, useSocks, proccessEmails, downloadEmails bool, maxMessagesToGet uint32, baseDir string) *Imaper {
	return &Imaper{
		SendersToSave:    sendersToSave,
		UseSocks:         useSocks,
		ProccessEmails:   proccessEmails,
		DownloadEmails:   downloadEmails,
		MaxMessagesToGet: maxMessagesToGet,
		MatcherBaseDir:   baseDir,
	}
}

// IMAPutil --
func (im *Imaper) IMAPutil(socksAddr string, addr string, emailUser string, emailPassword string) (bool, error) {
	var c *client.Client

	if im.UseSocks {
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

	if im.ProccessEmails {
		inboxProcessing(c, im, emailUser, emailPassword)
	}

	return true, nil
}

func inboxProcessing(c *client.Client, im *Imaper, emailUser, emailPassword string) {
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

	findMatchingSenders(mailboxList, c, im, emailUser, emailPassword)
}

func findMatchingSenders(mailboxList []string, c *client.Client, im *Imaper, emailUser, emailPassword string) {
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

		if mbox.Messages > im.MaxMessagesToGet {
			to = from + im.MaxMessagesToGet - 1
		}

		seqset := new(imap.SeqSet)
		seqset.AddRange(from, to)

		// Get the whole message body
		var section imap.BodySectionName
		items := []imap.FetchItem{section.FetchItem()}

		messages := make(chan *imap.Message, im.MaxMessagesToGet+1)
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

			for i := range im.SendersToSave {
				if strings.Contains(from[0].Address, im.SendersToSave[i]) {
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

						if !im.DownloadEmails {
							appendEmailCredentialsToFile(im.MatcherBaseDir, im.SendersToSave[i], emailUser, emailPassword)
							continue
						}

						switch h := p.Header.(type) {
						case *mail.InlineHeader:
							// This is the message's text (can be plain-text or HTML)
							b, err := ioutil.ReadAll(p.Body)
							if err != nil {
								log.Println("read body", err)
								continue
							}

							bb := strip.StripTags(string(b))
							//singleSpacePattern := regexp.MustCompile(`\s+`)
							//bbb := singleSpacePattern.ReplaceAllString(bb, " ")
							//log.Println(bbb)

							appendEmailBodyToFile(im.MatcherBaseDir, im.SendersToSave[i], bb, emailUser, emailPassword)
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

func appendEmailBodyToFile(basedir, fileName string, text, emailUser, emailPassword string) {
	var err error

	matcherDir := "matcherResults"

	err = utils.CheckDir(path.Join(basedir, matcherDir))
	utils.CheckError(err)

	filepath := path.Join(basedir, matcherDir, fileName+".txt")

	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	utils.CheckError(err)
	defer f.Close()

	// make email dividers
	dividerBegin := fmt.Sprintf("$#------------New Email from: %s:%s---------------#$", emailUser, emailPassword)
	dividerEnd := "$#------------End Email---------------#$"

	// wrap the email body with begin and end dividers
	textWithDividers := fmt.Sprintf("%s\r\n%s\r\n%s\r\n\r\n", dividerBegin, text, dividerEnd)

	_, err = f.WriteString(textWithDividers)
	utils.CheckError(err)
}

func appendEmailCredentialsToFile(basedir, fileName string, emailUser, emailPassword string) {
	var err error

	matcherDir := "matcherResults"

	err = utils.CheckDir(path.Join(basedir, matcherDir))
	utils.CheckError(err)

	filepath := path.Join(basedir, matcherDir, fileName+".txt")

	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	utils.CheckError(err)
	defer f.Close()

	line := fmt.Sprintf("%s:%s\n", emailUser, emailPassword)

	_, err = f.WriteString(line)
	utils.CheckError(err)
}
