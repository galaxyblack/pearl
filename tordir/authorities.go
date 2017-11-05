package tordir

import "github.com/erans/gonionoo"

// Reference: https://github.com/torproject/tor/blob/f755f9b9e67232d9d39682cbcdf4433ac738e17a/src/or/config.c#L1063-L1097
//
//	static const char *default_authorities[] = {
//	  "moria1 orport=9101 "
//	    "v3ident=D586D18309DED4CD6D57C18FDB97EFA96D330566 "
//	    "128.31.0.39:9131 9695 DFC3 5FFE B861 329B 9F1A B04C 4639 7020 CE31",
//	  "tor26 orport=443 "
//	    "v3ident=14C131DFC5C6F93646BE72FA1401C02A8DF2E8B4 "
//	    "ipv6=[2001:858:2:2:aabb:0:563b:1526]:443 "
//	    "86.59.21.38:80 847B 1F85 0344 D787 6491 A548 92F9 0493 4E4E B85D",
//	  "dizum orport=443 "
//	    "v3ident=E8A9C45EDE6D711294FADF8E7951F4DE6CA56B58 "
//	    "194.109.206.212:80 7EA6 EAD6 FD83 083C 538F 4403 8BBF A077 587D D755",
//	  "Bifroest orport=443 bridge "
//	    "37.218.247.217:80 1D8F 3A91 C37C 5D1C 4C19 B1AD 1D0C FBE8 BF72 D8E1",
//	  "gabelmoo orport=443 "
//	    "v3ident=ED03BB616EB2F60BEC80151114BB25CEF515B226 "
//	    "ipv6=[2001:638:a000:4140::ffff:189]:443 "
//	    "131.188.40.189:80 F204 4413 DAC2 E02E 3D6B CF47 35A1 9BCA 1DE9 7281",
//	  "dannenberg orport=443 "
//	    "v3ident=0232AF901C31A04EE9848595AF9BB7620D4C5B2E "
//	    "193.23.244.244:80 7BE6 83E6 5D48 1413 21C5 ED92 F075 C553 64AC 7123",
//	  "maatuska orport=80 "
//	    "v3ident=49015F787433103580E3B66A1707A00E60F2D15B "
//	    "ipv6=[2001:67c:289c::9]:80 "
//	    "171.25.193.9:443 BD6A 8292 55CB 08E6 6FBE 7D37 4836 3586 E46B 3810",
//	  "Faravahar orport=443 "
//	    "v3ident=EFCBE720AB3A82B99F9E953CD5BF50F7EEFC7B97 "
//	    "154.35.175.225:80 CF6D 0AAF B385 BE71 B8E1 11FC 5CFF 4B47 9237 33BC",
//	  "longclaw orport=443 "
//	    "v3ident=23D15D965BC35114467363C165C4F724B64B4F66 "
//	    "199.58.81.140:80 74A9 1064 6BCE EFBC D2E8 74FC 1DC9 9743 0F96 8145",
//	  "bastet orport=443 "
//	    "v3ident=27102BC123E7AF1D4741AE047E160C91ADC76B21 "
//	    "204.13.164.118:80 24E2 F139 121D 4394 C54B 5BCC 368B 3B41 1857 C413",
//	  NULL
//	};
//

// Authorities is a list of the directory addresses for the Tor directory
// authorities. This is unlikely to change often, but can be queried with the
// SearchAuthorityDirectoryAddresses() function. Listed at
// https://atlas.torproject.org/#search/flag:authority.
var Authorities = []string{
	"193.23.244.244:80",  // dannenberg
	"199.58.81.140:80",   // longclaw
	"194.109.206.212:80", // dizum
	"131.188.40.189:80",  // gabelmoo
	"86.59.21.38:80",     // tor26
	"37.218.247.217:80",  // Bifroest
	"154.35.175.225:80",  // Faravahar
	"128.31.0.34:9131",   // moria1
	"171.25.193.9:443",   // maatuska
	"204.13.164.118:80",  // bastet
}

// SearchAuthorityDirectoryAddresses queries the onionoo API for the directory
// addresses of the Tor authorities.
func SearchAuthorityDirectoryAddresses() ([]string, error) {
	query := map[string]string{
		"flag": "authority",
	}

	details, err := gonionoo.GetDetails(query)
	if err != nil {
		return nil, err
	}

	addresses := make([]string, len(details.Relays))
	for i, relay := range details.Relays {
		addresses[i] = relay.DirAddress
	}

	return addresses, nil
}
