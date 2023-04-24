package entities

type TYPE string
type VULNERABILITY string

const (
	Injection        TYPE = "INJECTION"
	Broken           TYPE = "BROKEN ACCESS CONTROL"
	Cryptography     TYPE = "CRYPTOGRAPHIC FAILURES"
	MisConfiguration TYPE = "SECURITY MISCONFIGURATION"
)

const (
	// sql injection
	SQLIboolean VULNERABILITY = "Boolean Based SQL Injection" // critical

	// xss
	CrossSiteScripting VULNERABILITY = "Cross-site Scripting" // hihg

	// broken access control
	LocalFileIncusion VULNERABILITY = "Local File Inclusion"       // low
	CrossSiteForgery  VULNERABILITY = "Cross-site Request Forgery" // low

	// cyptographic failure
	Certification VULNERABILITY = "SSL/TLS Not Implemented"             // medium
	Transmittion  VULNERABILITY = "Personal data Transmitted over HTTP" // medium

	// security misconfiguration
	PHPVersion   VULNERABILITY = "Version Disclosure (PHP)"   // low
	NginxVersion VULNERABILITY = "Version Disclosure (Nginx)" // low

	//
)

// PHP 5.3.0 to 7.0.0 vulnerabilities
const (
	PHPMemoryBuffer           VULNERABILITY = "PHP Improper Restriction of Operations within the Bounds of a Memory Buffer Vulnerability (CVE-2019-9638, CVE-2019-9639, CVE-2019-9641, CVE-2016-7480)" //
	PHPValidation             VULNERABILITY = "PHP Improper Input Validation Vulnerability (CVE-2017-8923)"                                                                                            //
	PHPNumericErrors          VULNERABILITY = "PHP Numeric Errors Vulnerability (CVE-2016-4344, CVE-2016-4345, CVE-2016-4346)"                                                                         //
	PHPIntegerOrWraparound    VULNERABILITY = "PHP Integer Overflow or Wraparound Vulnerability (CVE-2016-3078)"                                                                                       //
	PHPAccessControls         VULNERABILITY = "PHP Permissions, Privileges, and Access Controls Vulnerability (CVE-2019-9637)"                                                                         // low
	PHPAccessControl          VULNERABILITY = "PHP Improper Access Control Vulnerability (CVE-2016-5385)"                                                                                              // low
	PHPNullPointerDereference VULNERABILITY = "PHP NULL Pointer Dereference Vulnerability (CVE-2018-19395)"                                                                                            // low
	PHPDeserialization        VULNERABILITY = "PHP Deserialization of Untrusted Data Vulnerability (CVE-2018-19396)"
	PHPSSRF                   VULNERABILITY = "PHP Server-Side Request Forgery (SSRF) Vulnerability (CVE-2017-7272)"
	PHPThrottling             VULNERABILITY = "PHP Allocation of Resources Without Limits or Throttling Vulnerability (CVE-2017-7963)"
	PHPResourceConsumption    VULNERABILITY = "PHP Uncontrolled Resource Consumption Vulnerability (CVE-2015-9253)" // low
)
