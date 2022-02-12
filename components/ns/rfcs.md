# URN

<https://www.ietf.org/rfc/rfc2141.txt>

All URNs have the following syntax (phrases enclosed in quotes are REQUIRED):
<URN> ::= "urn:" <NID> ":" <NSS>
where <NID> is the Namespace Identifier, and <NSS> is the Namespace Specific String.  The leading "urn:" sequence
is case-insensitive. The Namespace ID determines the _syntactic_ interpretation of the Namespace Specific String

The following URN comparisons highlight the lexical equivalence definitions:

    URN:foo:a123,456
    urn:foo:a123,456
    urn:FOO:a123,456
    urn:foo:A123,456
    urn:foo:a123%2C456
    URN:FOO:a123%2c456

## Golang implementation

<https://github.com/voicera/gooseberry/urn>

    type URN struct {
        // contains unexported fields
    }

    func NewURN(namespaceID string, namespaceSpecificString string) *URN
    func TryParseString(urn string) (*URN, bool)
    func (urn *URN) GetNamespaceID() string
    func (urn *URN) GetNamespaceSpecificString() string
    func (urn *URN) MarshalJSON() ([]byte, error)
    func (urn *URN) String() string
    func (urn *URN) UnmarshalJSON(b []byte) error


# URI, URL 

<https://www.ietf.org/rfc/rfc3986.txt>     

A URI can be further classified as a locator, a name, or both. The term "Uniform Resource Locator" (URL) refers to the
subset of URIs that, in addition to identifying a resource, provide a means of locating the resource by describing its
primary access mechanism (e.g., its network "location"). The term "Uniform Resource Name" (URN) has been used historically
to refer to both URIs under the "urn" scheme [RFC2141], which are required to remain globally unique and persistent
even when the resource ceases to exist or becomes unavailable, and to any other URI with the properties of a name.

Examples:

    ftp://ftp.is.co.za/rfc/rfc1808.txt
    http://www.ietf.org/rfc/rfc2396.txt
    ldap://[2001:db8::7]/c=GB?objectClass?one
    mailto:John.Doe@example.com
    news:comp.infosystems.www.servers.unix
    tel:+1-816-555-1212
    telnet://192.0.2.16:80/
    urn:oasis:names:specification:docbook:dtd:xml:4.1.2

## Golang implementation (url.URL)

    type URL struct {
        Scheme      string
        Opaque      string    // encoded opaque data
        User        *Userinfo // username and password information
        Host        string    // host or host:port
        Path        string    // path (relative paths may omit leading slash)
        RawPath     string    // encoded path hint (see EscapedPath method)
        ForceQuery  bool      // append a query ('?') even if RawQuery is empty
        RawQuery    string    // encoded query values, without '?'
        Fragment    string    // fragment for references, without '#'
        RawFragment string    // encoded fragment hint (see EscapedFragment method)
    }

