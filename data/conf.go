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
	VALIDFILE       Value
	NOTFOUNDFILE    Value
	HOSTFILE        Value
	MATCHERFILE     Value
	SOCKSFILE       Value
	MAXJOBS         Value
	BUFFERSIZE      Value
	SAVELASTLINELOG Value
}

// NewUserValues return new parser object
func NewUserValues() *UserValues {
	return &UserValues{}
}

// Open the json file
func (uv *UserValues) Open(filename string) error {
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

// Write the json file
func (uv *UserValues) Write(filename string) error {
	out, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer out.Close()

	encodeJSON := json.NewEncoder(out)
	err = encodeJSON.Encode(uv)
	if err != nil {
		return err
	}

	return nil
}

// GetVALIDFILE returns the filename of VALID MAILPASS file
func (uv *UserValues) GetVALIDFILE() string {
	return uv.VALIDFILE.String()
}

// SetVALIDFILE sets a new filename for the VALID MAILPASS file
func (uv *UserValues) SetVALIDFILE(filename string) {
	uv.VALIDFILE = Value(filename)
}

// GetNOTFOUNDFILE returns the filename of NOTFOUND MAILPASS file
func (uv *UserValues) GetNOTFOUNDFILE() string {
	return uv.NOTFOUNDFILE.String()
}

// SetNOTFOUNDFILE sets a new filename for the NOTFOUNDFILE MAILPASS file
func (uv *UserValues) SetNOTFOUNDFILE(filename string) {
	uv.NOTFOUNDFILE = Value(filename)
}

// GetHOSTFILE returns the filename of Hostfile
func (uv *UserValues) GetHOSTFILE() string {
	return uv.HOSTFILE.String()
}

// SetHOSTFILE sets a new filename for the Hostfile
func (uv *UserValues) SetHOSTFILE(filename string) {
	uv.HOSTFILE = Value(filename)
}

// GetMATCHERFILE returns the filename of the Matcherfile
func (uv *UserValues) GetMATCHERFILE() string {
	return uv.MATCHERFILE.String()
}

// SetMATCHERFILE sets a new filename for the Matcherfile
func (uv *UserValues) SetMATCHERFILE(filename string) {
	uv.MATCHERFILE = Value(filename)
}

// GetSOCKSFILE returns the filename of the Socksfile
func (uv *UserValues) GetSOCKSFILE() string {
	return uv.SOCKSFILE.String()
}

// SetSOCKSFILE sets a new filename for the Socksfile
func (uv *UserValues) SetSOCKSFILE(filename string) {
	uv.SOCKSFILE = Value(filename)
}

// GetMAXJOBS returns the num of max jobs hold in memory
func (uv *UserValues) GetMAXJOBS() int {
	return uv.MAXJOBS.Int()
}

// SetMAXJOBS sets num of max jobs hold in memory
func (uv *UserValues) SetMAXJOBS(n int) {
	uv.MAXJOBS = Value(strconv.Itoa(n))
}

// GetBUFFERSIZE returns the amount of bytes holding in the memory
func (uv *UserValues) GetBUFFERSIZE() int {
	return uv.BUFFERSIZE.Int()
}

// SetBUFFERSIZE sets the amount of bytes holding in the memory
func (uv *UserValues) SetBUFFERSIZE(n int) {
	uv.BUFFERSIZE = Value(strconv.Itoa(n))
}

// IsSAVELASTLINELOG returns if LASTLINELOG will be saved
func (uv *UserValues) IsSAVELASTLINELOG() bool {
	return uv.SAVELASTLINELOG.ToBool()
}

// SetSAVELASTLINELOG set the bool value if a lastlinelog is saved or not
func (uv *UserValues) SetSAVELASTLINELOG(b bool) {
	uv.SAVELASTLINELOG = Value(strconv.FormatBool(b))
}

//----------------- General Config --------------------

// Config is our general config struct and holds the parsed UserValues
type Config struct {
	WORKERS      Value
	USESOCKS     bool
	PROCESSMAILS bool
	USERVALUE    *UserValues
}

// NewConf return new parser object
func NewConf() *Config {
	return &Config{}
}

// Open parses the JSON UserValues and initialize a new general config.
// Also NewConf set GOMAXPROCS to our number of available CPU cores on this machine
// It
func (c *Config) Open(filename string) error {
	var err error

	// Now lets parse our UserValue json config file
	uval := NewUserValues()
	err = uval.Open(filename)
	if err != nil {
		return err
	}

	// Now lets find out how many workers we can effectivly use on this computer
	// Set GOMAXPROCS on this value too
	cores := runtime.NumCPU()
	_ = runtime.GOMAXPROCS(cores)

	// Set the values
	c.WORKERS = Value(strconv.Itoa(cores))
	c.USESOCKS = uval.GetSOCKSFILE() != ""
	c.PROCESSMAILS = uval.GetMATCHERFILE() != ""
	c.USERVALUE = uval

	return nil
}

// GetWorkers returns the currently num of cpu workers
func (c *Config) GetWorkers() int {
	return c.WORKERS.Int()
}

// SetWorkers set another num of workers as system default
// AWARE, this should not be used in production!!
func (c *Config) SetWorkers(workers int) {
	c.WORKERS = Value(strconv.Itoa(workers))
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
