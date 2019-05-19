package getcert

const (
	certDotPemPath = "/var/www/cert/fullchain.pem"
	keyDotPemPath  = "/var/www/cert/key.pem"
)

//DomainCert defines domain's  cert object
type DomainCert struct {
	domain       string
	certFilePath string
	keyFilePath  string
}

//NewDomainCert new domain object
func NewDomainCert(domainName string) *DomainCert {
	return &DomainCert{
		domain:       domainName,
		certFilePath: certDotPemPath,
		keyFilePath:  keyDotPemPath,
	}
}

//GetCert get cert file path
func (d *DomainCert) GetCert() string {
	return d.certFilePath
}

//GetKey get key.pem
func (d *DomainCert) GetKey() string {
	return d.keyFilePath
}

//GetDomain get cert domain which sets before
func (d *DomainCert) GetDomain() string {
	return d.domain
}
