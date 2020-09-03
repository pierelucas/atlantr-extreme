package data

import (
	"encoding/json"
	"os"
	"runtime"
	"strconv"
)

//----------------- User Config --------------------

// UserValues json key
type UserValues struct {
	MAILPASS     value
	VALIDFILE    value
	NOTFOUNDFILE value
	HOSTFILE     value
	MATCHERFILE  value
	SOCKSFILE    value
	MAXJOBS      value
	BUFFERSIZE   value
}

// NewUserValues return new parser object
func NewUserValues() *UserValues {
	return &UserValues{}
}

// Parse the json file
func (uv *UserValues) Parse(filename string) error {
	in, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer in.Close()

	decodeJSON := json.NewDecoder(in)
	err = decodeJSON.Decode(uv)
	if err != nil {
		return err
	}

	return nil
}

// GetMAILPASS returns the filename of MAILPASS file
func (uv *UserValues) GetMAILPASS() string {
	return uv.MAILPASS.String()
}

// SetMAILPASS sets a new filename for the MAILPASS file
func (uv *UserValues) SetMAILPASS(filename string) {
	uv.MAILPASS = value(filename)
}

// GetVALIDFILE returns the filename of VALID MAILPASS file
func (uv *UserValues) GetVALIDFILE() string {
	return uv.VALIDFILE.String()
}

// SetVALIDFILE sets a new filename for the VALID MAILPASS file
func (uv *UserValues) SetVALIDFILE(filename string) {
	uv.VALIDFILE = value(filename)
}

// GetNOTFOUNDFILE returns the filename of NOTFOUND MAILPASS file
func (uv *UserValues) GetNOTFOUNDFILE() string {
	return uv.NOTFOUNDFILE.String()
}

// SetNOTFOUNDFILE sets a new filename for the NOTFOUNDFILE MAILPASS file
func (uv *UserValues) SetNOTFOUNDFILE(filename string) {
	uv.NOTFOUNDFILE = value(filename)
}

// GetHOSTFILE returns the filename of Hostfile
func (uv *UserValues) GetHOSTFILE() string {
	return uv.HOSTFILE.String()
}

// SetHOSTFILE sets a new filename for the Hostfile
func (uv *UserValues) SetHOSTFILE(filename string) {
	uv.HOSTFILE = value(filename)
}

// GetMATCHERFILE returns the filename of the Matcherfile
func (uv *UserValues) GetMATCHERFILE() string {
	return uv.MATCHERFILE.String()
}

// SetMATCHERFILE sets a new filename for the Matcherfile
func (uv *UserValues) SetMATCHERFILE(filename string) {
	uv.MATCHERFILE = value(filename)
}

// GetSOCKSFILE returns the filename of the Socksfile
func (uv *UserValues) GetSOCKSFILE() string {
	return uv.SOCKSFILE.String()
}

// SetSOCKSFILE sets a new filename for the Socksfile
func (uv *UserValues) SetSOCKSFILE(filename string) {
	uv.SOCKSFILE = value(filename)
}

// GetMAXJOBS returns the num of max jobs hold in memory
func (uv *UserValues) GetMAXJOBS() int {
	return uv.MAXJOBS.Int()
}

// SetMAXJOBS sets num of max jobs hold in memory
func (uv *UserValues) SetMAXJOBS(n int) {
	uv.MAXJOBS = value(strconv.Itoa(n))
}

// GetBUFFERSIZE returns the amount of bytes holding in the memory
func (uv *UserValues) GetBUFFERSIZE() int {
	return uv.BUFFERSIZE.Int()
}

// SetBUFFERSIZE sets the amount of bytes holding in the memory
func (uv *UserValues) SetBUFFERSIZE(n int) {
	uv.BUFFERSIZE = value(strconv.Itoa(n))
}

//----------------- General Config --------------------

// Config is our general config struct and holds the parsed UserValues
type Config struct {
	WORKERS      value
	USESOCKS     bool
	PROCESSMAILS bool
	USERVALUE    *UserValues
}

// NewConf parses the JSON UserValues and initialize a new general config.
// Also NewConf set GOMAXPROCS to our number of available CPU cores on this machine
// It
func NewConf(filename string) (*Config, error) {
	var err error

	// Now lets parse our UserValue json config file
	uval := NewUserValues()
	err = uval.Parse(filename)
	if err != nil {
		return nil, err
	}

	// Now lets find out how many workers we can effectivly use on this computer
	// Set GOMAXPROCS on this value too
	cores := runtime.NumCPU()
	_ = runtime.GOMAXPROCS(cores)

	return &Config{
		WORKERS:      value(strconv.Itoa(cores)),
		USESOCKS:     uval.GetSOCKSFILE() != "",
		PROCESSMAILS: uval.GetMATCHERFILE() != "",
		USERVALUE:    uval,
	}, nil
}

// GetWorkers returns the currently num of cpu workers
func (c *Config) GetWorkers() int {
	return c.WORKERS.Int()
}

// SetWorkers set another num of workers as system default
// AWARE, this should not be used in production!!
func (c *Config) SetWorkers(workers int) {
	c.WORKERS = value(strconv.Itoa(workers))
}

// GetUSESOCKS returns if a socks proxy is used or not
func (c *Config) GetUSESOCKS() bool {
	return c.USESOCKS
}

// SetUSESOCKS sets if we use a socks proxy or not
func (c *Config) SetUSESOCKS(b bool) {
	c.USESOCKS = b
}

// GetPROCESSMAILS returns if we process the inbox mail with matchers or not
func (c *Config) GetPROCESSMAILS() bool {
	return c.PROCESSMAILS
}

// SetPROCESSMAILS sets if we process the inbox mail with mathcers or not
func (c *Config) SetPROCESSMAILS(b bool) {
	c.PROCESSMAILS = b
}
