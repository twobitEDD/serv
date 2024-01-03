// This accounts file contains the accounts and corresponding token allocation for
// accounts which participated in the Olympus Mons Testnet (November 2021) through
// completion of the Mars Meteor Missions. The token allocation is in aserv.

// 7.5% of the genesis allocation, totaling to ~7.5M Serv tokens, was set aside for
// participants of the incentivized testnet, of which ~5.6M is claimed here. The
// remaining funds will be sent to the community pool.

// The participants of the testnet will immediately stake their tokens to a set
// of chosen validators, also included in this file.

package v11

// Allocations are rewards by participant address, along with the randomized validator to stake to.
var Allocations = [1][3]string{
	{
		"sx1az0eeqw5222en0enxz9ect3s8xy4kayz7khcn2",
		"72232472320000001048576",
		"evmosvaloper1aep37tvd5yh4ydeya2a88g2tkjz6f2lcrfjumh",
	},
}
