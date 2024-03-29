package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	"github.com/twobitedd/serv/v12/app"
	ibctesting "github.com/twobitedd/serv/v12/ibc/testing"
	"github.com/twobitedd/serv/v12/testutil"
	teststypes "github.com/twobitedd/serv/v12/types/tests"
	claimstypes "github.com/twobitedd/serv/v12/x/claims/types"
	"github.com/twobitedd/serv/v12/x/recovery/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Recovery: Performing an IBC Transfer", Ordered, func() {
	coinServ := sdk.NewCoin("aserv", sdk.NewInt(10000))
	coinOsmo := sdk.NewCoin("uosmo", sdk.NewInt(10))
	coinAtom := sdk.NewCoin("uatom", sdk.NewInt(10))

	var (
		sender, receiver       string
		senderAcc, receiverAcc sdk.AccAddress
		timeout                uint64
		claim                  claimstypes.ClaimsRecord
	)

	BeforeEach(func() {
		s.SetupTest()
	})

	Describe("from a non-authorized chain", func() {
		BeforeEach(func() {
			params := claimstypes.DefaultParams()
			params.AuthorizedChannels = []string{}
			err := s.ServChain.App.(*app.Serv).ClaimsKeeper.SetParams(s.ServChain.GetContext(), params)
			s.Require().NoError(err)

			sender = s.IBCOsmosisChain.SenderAccount.GetAddress().String()
			receiver = s.ServChain.SenderAccount.GetAddress().String()
			senderAcc = sdk.MustAccAddressFromBech32(sender)
			receiverAcc = sdk.MustAccAddressFromBech32(receiver)
		})
		It("should transfer and not recover tokens", func() {
			s.SendAndReceiveMessage(s.pathOsmosisServ, s.IBCOsmosisChain, "uosmo", 10, sender, receiver, 1)

			nativeServ := s.ServChain.App.(*app.Serv).BankKeeper.GetBalance(s.ServChain.GetContext(), senderAcc, "aserv")
			Expect(nativeServ).To(Equal(coinServ))
			ibcOsmo := s.ServChain.App.(*app.Serv).BankKeeper.GetBalance(s.ServChain.GetContext(), receiverAcc, teststypes.UosmoIbcdenom)
			Expect(ibcOsmo).To(Equal(sdk.NewCoin(teststypes.UosmoIbcdenom, coinOsmo.Amount)))
		})
	})

	Describe("from an authorized, non-EVM chain (e.g. Osmosis)", func() {
		Describe("to a different account on Serv (sender != recipient)", func() {
			BeforeEach(func() {
				sender = s.IBCOsmosisChain.SenderAccount.GetAddress().String()
				receiver = s.ServChain.SenderAccount.GetAddress().String()
				senderAcc = sdk.MustAccAddressFromBech32(sender)
				receiverAcc = sdk.MustAccAddressFromBech32(receiver)
			})

			It("should transfer and not recover tokens", func() {
				s.SendAndReceiveMessage(s.pathOsmosisServ, s.IBCOsmosisChain, "uosmo", 10, sender, receiver, 1)

				nativeServ := s.ServChain.App.(*app.Serv).BankKeeper.GetBalance(s.ServChain.GetContext(), senderAcc, "aserv")
				Expect(nativeServ).To(Equal(coinServ))
				ibcOsmo := s.ServChain.App.(*app.Serv).BankKeeper.GetBalance(s.ServChain.GetContext(), receiverAcc, teststypes.UosmoIbcdenom)
				Expect(ibcOsmo).To(Equal(sdk.NewCoin(teststypes.UosmoIbcdenom, coinOsmo.Amount)))
			})
		})

		Describe("to the sender's own eth_secp256k1 account on Serv (sender == recipient)", func() {
			BeforeEach(func() {
				sender = s.IBCOsmosisChain.SenderAccount.GetAddress().String()
				receiver = s.IBCOsmosisChain.SenderAccount.GetAddress().String()
				senderAcc = sdk.MustAccAddressFromBech32(sender)
				receiverAcc = sdk.MustAccAddressFromBech32(receiver)
			})

			Context("with disabled recovery parameter", func() {
				BeforeEach(func() {
					params := types.DefaultParams()
					params.EnableRecovery = false
					s.ServChain.App.(*app.Serv).RecoveryKeeper.SetParams(s.ServChain.GetContext(), params) //nolint:errcheck
				})

				It("should not transfer or recover tokens", func() {
					s.SendAndReceiveMessage(s.pathOsmosisServ, s.IBCOsmosisChain, coinOsmo.Denom, coinOsmo.Amount.Int64(), sender, receiver, 1)
					nativeServ := s.ServChain.App.(*app.Serv).BankKeeper.GetBalance(s.ServChain.GetContext(), senderAcc, "aserv")
					Expect(nativeServ).To(Equal(coinServ))
					ibcOsmo := s.ServChain.App.(*app.Serv).BankKeeper.GetBalance(s.ServChain.GetContext(), receiverAcc, teststypes.UosmoIbcdenom)
					Expect(ibcOsmo).To(Equal(sdk.NewCoin(teststypes.UosmoIbcdenom, coinOsmo.Amount)))
				})
			})

			Context("with a sender's claims record", func() {
				Context("without completed actions", func() {
					BeforeEach(func() {
						amt := sdk.NewInt(int64(100))
						coins := sdk.NewCoins(sdk.NewCoin("aserv", amt))
						claim = claimstypes.NewClaimsRecord(amt)
						s.ServChain.App.(*app.Serv).ClaimsKeeper.SetClaimsRecord(s.ServChain.GetContext(), senderAcc, claim)

						// update the escrowed account balance to maintain the invariant
						err := testutil.FundModuleAccount(s.ServChain.GetContext(), s.ServChain.App.(*app.Serv).BankKeeper, claimstypes.ModuleName, coins)
						s.Require().NoError(err)
					})

					It("should not transfer or recover tokens", func() {
						// Prevent further funds from getting stuck
						s.SendAndReceiveMessage(s.pathOsmosisServ, s.IBCOsmosisChain, coinOsmo.Denom, coinOsmo.Amount.Int64(), sender, receiver, 1)
						nativeServ := s.ServChain.App.(*app.Serv).BankKeeper.GetBalance(s.ServChain.GetContext(), senderAcc, "aserv")
						Expect(nativeServ).To(Equal(coinServ))
						ibcOsmo := s.ServChain.App.(*app.Serv).BankKeeper.GetBalance(s.ServChain.GetContext(), receiverAcc, teststypes.UosmoIbcdenom)
						Expect(ibcOsmo.IsZero()).To(BeTrue())
					})
				})

				Context("with completed actions", func() {
					// Already has stuck funds
					BeforeEach(func() {
						amt := sdk.NewInt(int64(100))
						coins := sdk.NewCoins(sdk.NewCoin("aserv", sdk.NewInt(int64(75))))
						claim = claimstypes.NewClaimsRecord(amt)
						claim.MarkClaimed(claimstypes.ActionIBCTransfer)
						s.ServChain.App.(*app.Serv).ClaimsKeeper.SetClaimsRecord(s.ServChain.GetContext(), senderAcc, claim)

						// update the escrowed account balance to maintain the invariant
						err := testutil.FundModuleAccount(s.ServChain.GetContext(), s.ServChain.App.(*app.Serv).BankKeeper, claimstypes.ModuleName, coins)
						s.Require().NoError(err)

						// aserv & ibc tokens that originated from the sender's chain
						s.SendAndReceiveMessage(s.pathOsmosisServ, s.IBCOsmosisChain, coinOsmo.Denom, coinOsmo.Amount.Int64(), sender, receiver, 1)
						timeout = uint64(s.ServChain.GetContext().BlockTime().Add(time.Hour * 4).Add(time.Second * -20).UnixNano())
					})

					It("should transfer tokens to the recipient and perform recovery", func() {
						// Escrow before relaying packets
						balanceEscrow := s.ServChain.App.(*app.Serv).BankKeeper.GetBalance(s.ServChain.GetContext(), transfertypes.GetEscrowAddress("transfer", "channel-0"), "aserv")
						Expect(balanceEscrow).To(Equal(coinServ))
						ibcOsmo := s.ServChain.App.(*app.Serv).BankKeeper.GetBalance(s.ServChain.GetContext(), receiverAcc, teststypes.UosmoIbcdenom)
						Expect(ibcOsmo.IsZero()).To(BeTrue())

						// Relay both packets that were sent in the ibc_callback
						err := s.pathOsmosisServ.RelayPacket(CreatePacket("10000", "aserv", sender, receiver, "transfer", "channel-0", "transfer", "channel-0", 1, timeout))
						s.Require().NoError(err)
						err = s.pathOsmosisServ.RelayPacket(CreatePacket("10", "transfer/channel-0/uosmo", sender, receiver, "transfer", "channel-0", "transfer", "channel-0", 2, timeout))
						s.Require().NoError(err)

						// Check that the aserv were recovered
						nativeServ := s.ServChain.App.(*app.Serv).BankKeeper.GetBalance(s.ServChain.GetContext(), senderAcc, "aserv")
						Expect(nativeServ.IsZero()).To(BeTrue())
						ibcServ := s.IBCOsmosisChain.GetSimApp().BankKeeper.GetBalance(s.IBCOsmosisChain.GetContext(), receiverAcc, teststypes.AservIbcdenom)
						Expect(ibcServ).To(Equal(sdk.NewCoin(teststypes.AservIbcdenom, coinServ.Amount)))

						// Check that the uosmo were recovered
						ibcOsmo = s.ServChain.App.(*app.Serv).BankKeeper.GetBalance(s.ServChain.GetContext(), receiverAcc, teststypes.UosmoIbcdenom)
						Expect(ibcOsmo.IsZero()).To(BeTrue())
						nativeOsmo := s.IBCOsmosisChain.GetSimApp().BankKeeper.GetBalance(s.IBCOsmosisChain.GetContext(), receiverAcc, "uosmo")
						Expect(nativeOsmo).To(Equal(coinOsmo))
					})

					It("should not claim/migrate/merge claims records", func() {
						// Relay both packets that were sent in the ibc_callback
						err := s.pathOsmosisServ.RelayPacket(CreatePacket("10000", "aserv", sender, receiver, "transfer", "channel-0", "transfer", "channel-0", 1, timeout))
						s.Require().NoError(err)
						err = s.pathOsmosisServ.RelayPacket(CreatePacket("10", "transfer/channel-0/uosmo", sender, receiver, "transfer", "channel-0", "transfer", "channel-0", 2, timeout))
						s.Require().NoError(err)

						claimAfter, _ := s.ServChain.App.(*app.Serv).ClaimsKeeper.GetClaimsRecord(s.ServChain.GetContext(), senderAcc)
						Expect(claim).To(Equal(claimAfter))
					})
				})
			})

			Context("without a sender's claims record", func() {
				When("recipient has no ibc vouchers that originated from other chains", func() {
					It("should transfer and recover tokens", func() {
						// aserv & ibc tokens that originated from the sender's chain
						s.SendAndReceiveMessage(s.pathOsmosisServ, s.IBCOsmosisChain, coinOsmo.Denom, coinOsmo.Amount.Int64(), sender, receiver, 1)
						timeout = uint64(s.ServChain.GetContext().BlockTime().Add(time.Hour * 4).Add(time.Second * -20).UnixNano())

						// Escrow before relaying packets
						balanceEscrow := s.ServChain.App.(*app.Serv).BankKeeper.GetBalance(s.ServChain.GetContext(), transfertypes.GetEscrowAddress("transfer", "channel-0"), "aserv")
						Expect(balanceEscrow).To(Equal(coinServ))
						ibcOsmo := s.ServChain.App.(*app.Serv).BankKeeper.GetBalance(s.ServChain.GetContext(), receiverAcc, teststypes.UosmoIbcdenom)
						Expect(ibcOsmo.IsZero()).To(BeTrue())

						// Relay both packets that were sent in the ibc_callback
						err := s.pathOsmosisServ.RelayPacket(CreatePacket("10000", "aserv", sender, receiver, "transfer", "channel-0", "transfer", "channel-0", 1, timeout))
						s.Require().NoError(err)
						err = s.pathOsmosisServ.RelayPacket(CreatePacket("10", "transfer/channel-0/uosmo", sender, receiver, "transfer", "channel-0", "transfer", "channel-0", 2, timeout))
						s.Require().NoError(err)

						// Check that the aserv were recovered
						nativeServ := s.ServChain.App.(*app.Serv).BankKeeper.GetBalance(s.ServChain.GetContext(), senderAcc, "aserv")
						Expect(nativeServ.IsZero()).To(BeTrue())
						ibcServ := s.IBCOsmosisChain.GetSimApp().BankKeeper.GetBalance(s.IBCOsmosisChain.GetContext(), receiverAcc, teststypes.AservIbcdenom)
						Expect(ibcServ).To(Equal(sdk.NewCoin(teststypes.AservIbcdenom, coinServ.Amount)))

						// Check that the uosmo were recovered
						ibcOsmo = s.ServChain.App.(*app.Serv).BankKeeper.GetBalance(s.ServChain.GetContext(), receiverAcc, teststypes.UosmoIbcdenom)
						Expect(ibcOsmo.IsZero()).To(BeTrue())
						nativeOsmo := s.IBCOsmosisChain.GetSimApp().BankKeeper.GetBalance(s.IBCOsmosisChain.GetContext(), receiverAcc, "uosmo")
						Expect(nativeOsmo).To(Equal(coinOsmo))
					})
				})

				// Do not recover uatom sent from Cosmos when performing recovery through IBC transfer from Osmosis
				When("recipient has additional ibc vouchers that originated from other chains", func() {
					BeforeEach(func() {
						params := types.DefaultParams()
						params.EnableRecovery = false
						s.ServChain.App.(*app.Serv).RecoveryKeeper.SetParams(s.ServChain.GetContext(), params) //nolint:errcheck

						// Send uatom from Cosmos to Serv
						s.SendAndReceiveMessage(s.pathCosmosServ, s.IBCCosmosChain, coinAtom.Denom, coinAtom.Amount.Int64(), s.IBCCosmosChain.SenderAccount.GetAddress().String(), receiver, 1)

						params.EnableRecovery = true
						s.ServChain.App.(*app.Serv).RecoveryKeeper.SetParams(s.ServChain.GetContext(), params) //nolint:errcheck
					})
					It("should not recover tokens that originated from other chains", func() {
						// Send uosmo from Osmosis to Serv
						s.SendAndReceiveMessage(s.pathOsmosisServ, s.IBCOsmosisChain, "uosmo", 10, sender, receiver, 1)

						// Relay both packets that were sent in the ibc_callback
						timeout := uint64(s.ServChain.GetContext().BlockTime().Add(time.Hour * 4).Add(time.Second * -20).UnixNano())
						err := s.pathOsmosisServ.RelayPacket(CreatePacket("10000", "aserv", sender, receiver, "transfer", "channel-0", "transfer", "channel-0", 1, timeout))
						s.Require().NoError(err)
						err = s.pathOsmosisServ.RelayPacket(CreatePacket("10", "transfer/channel-0/uosmo", sender, receiver, "transfer", "channel-0", "transfer", "channel-0", 2, timeout))
						s.Require().NoError(err)

						// Aserv was recovered from user address
						nativeServ := s.ServChain.App.(*app.Serv).BankKeeper.GetBalance(s.ServChain.GetContext(), senderAcc, "aserv")
						Expect(nativeServ.IsZero()).To(BeTrue())
						ibcServ := s.IBCOsmosisChain.GetSimApp().BankKeeper.GetBalance(s.IBCOsmosisChain.GetContext(), receiverAcc, teststypes.AservIbcdenom)
						Expect(ibcServ).To(Equal(sdk.NewCoin(teststypes.AservIbcdenom, coinServ.Amount)))

						// Check that the uosmo were retrieved
						ibcOsmo := s.ServChain.App.(*app.Serv).BankKeeper.GetBalance(s.ServChain.GetContext(), receiverAcc, teststypes.UosmoIbcdenom)
						Expect(ibcOsmo.IsZero()).To(BeTrue())
						nativeOsmo := s.IBCOsmosisChain.GetSimApp().BankKeeper.GetBalance(s.IBCOsmosisChain.GetContext(), receiverAcc, "uosmo")
						Expect(nativeOsmo).To(Equal(coinOsmo))

						// Check that the atoms were not retrieved
						ibcAtom := s.ServChain.App.(*app.Serv).BankKeeper.GetBalance(s.ServChain.GetContext(), senderAcc, teststypes.UatomIbcdenom)
						Expect(ibcAtom).To(Equal(sdk.NewCoin(teststypes.UatomIbcdenom, coinAtom.Amount)))

						// Repeat transaction from Osmosis to Serv
						s.SendAndReceiveMessage(s.pathOsmosisServ, s.IBCOsmosisChain, "uosmo", 10, sender, receiver, 2)

						timeout = uint64(s.ServChain.GetContext().BlockTime().Add(time.Hour * 4).Add(time.Second * -20).UnixNano())
						err = s.pathOsmosisServ.RelayPacket(CreatePacket("10", "transfer/channel-0/uosmo", sender, receiver, "transfer", "channel-0", "transfer", "channel-0", 3, timeout))
						s.Require().NoError(err)

						// No further tokens recovered
						nativeServ = s.ServChain.App.(*app.Serv).BankKeeper.GetBalance(s.ServChain.GetContext(), senderAcc, "aserv")
						Expect(nativeServ.IsZero()).To(BeTrue())
						ibcServ = s.IBCOsmosisChain.GetSimApp().BankKeeper.GetBalance(s.IBCOsmosisChain.GetContext(), receiverAcc, teststypes.AservIbcdenom)
						Expect(ibcServ).To(Equal(sdk.NewCoin(teststypes.AservIbcdenom, coinServ.Amount)))

						ibcOsmo = s.ServChain.App.(*app.Serv).BankKeeper.GetBalance(s.ServChain.GetContext(), receiverAcc, teststypes.UosmoIbcdenom)
						Expect(ibcOsmo.IsZero()).To(BeTrue())
						nativeOsmo = s.IBCOsmosisChain.GetSimApp().BankKeeper.GetBalance(s.IBCOsmosisChain.GetContext(), receiverAcc, "uosmo")
						Expect(nativeOsmo).To(Equal(coinOsmo))

						ibcAtom = s.ServChain.App.(*app.Serv).BankKeeper.GetBalance(s.ServChain.GetContext(), senderAcc, teststypes.UatomIbcdenom)
						Expect(ibcAtom).To(Equal(sdk.NewCoin(teststypes.UatomIbcdenom, coinAtom.Amount)))
					})
				})

				// Recover ibc/uatom that was sent from Osmosis back to Osmosis
				When("recipient has additional non-native ibc vouchers that originated from senders chains", func() {
					BeforeEach(func() {
						params := types.DefaultParams()
						params.EnableRecovery = false
						err := s.ServChain.App.(*app.Serv).RecoveryKeeper.SetParams(s.ServChain.GetContext(), params)
						s.Require().NoError(err)
						s.SendAndReceiveMessage(s.pathOsmosisCosmos, s.IBCCosmosChain, coinAtom.Denom, coinAtom.Amount.Int64(), s.IBCCosmosChain.SenderAccount.GetAddress().String(), receiver, 1)

						// Send IBC transaction of 10 ibc/uatom
						transferMsg := transfertypes.NewMsgTransfer(s.pathOsmosisServ.EndpointA.ChannelConfig.PortID, s.pathOsmosisServ.EndpointA.ChannelID, sdk.NewCoin(teststypes.UatomIbcdenom, sdk.NewInt(10)), sender, receiver, timeoutHeight, 0, "")
						_, err = ibctesting.SendMsgs(s.IBCOsmosisChain, ibctesting.DefaultFeeAmt, transferMsg)
						s.Require().NoError(err) // message committed
						transfer := transfertypes.NewFungibleTokenPacketData("transfer/channel-1/uatom", "10", sender, receiver, "")
						packet := channeltypes.NewPacket(transfer.GetBytes(), 1, s.pathOsmosisServ.EndpointA.ChannelConfig.PortID, s.pathOsmosisServ.EndpointA.ChannelID, s.pathOsmosisServ.EndpointB.ChannelConfig.PortID, s.pathOsmosisServ.EndpointB.ChannelID, timeoutHeight, 0)
						// Receive message on the serv side, and send ack
						err = s.pathOsmosisServ.RelayPacket(packet)
						s.Require().NoError(err)

						// Check that the ibc/uatom are available
						osmoIBCAtom := s.ServChain.App.(*app.Serv).BankKeeper.GetBalance(s.ServChain.GetContext(), receiverAcc, teststypes.UatomOsmoIbcdenom)
						s.Require().Equal(osmoIBCAtom.Amount, coinAtom.Amount)

						params.EnableRecovery = true
						s.ServChain.App.(*app.Serv).RecoveryKeeper.SetParams(s.ServChain.GetContext(), params) //nolint:errcheck
					})
					It("should not recover tokens that originated from other chains", func() {
						s.SendAndReceiveMessage(s.pathOsmosisServ, s.IBCOsmosisChain, "uosmo", 10, sender, receiver, 2)

						// Relay packets that were sent in the ibc_callback
						timeout := uint64(s.ServChain.GetContext().BlockTime().Add(time.Hour * 4).Add(time.Second * -20).UnixNano())
						err := s.pathOsmosisServ.RelayPacket(CreatePacket("10000", "aserv", sender, receiver, "transfer", "channel-0", "transfer", "channel-0", 1, timeout))
						s.Require().NoError(err)
						err = s.pathOsmosisServ.RelayPacket(CreatePacket("10", "transfer/channel-0/transfer/channel-1/uatom", sender, receiver, "transfer", "channel-0", "transfer", "channel-0", 2, timeout))
						s.Require().NoError(err)
						err = s.pathOsmosisServ.RelayPacket(CreatePacket("10", "transfer/channel-0/uosmo", sender, receiver, "transfer", "channel-0", "transfer", "channel-0", 3, timeout))
						s.Require().NoError(err)

						// Aserv was recovered from user address
						nativeServ := s.ServChain.App.(*app.Serv).BankKeeper.GetBalance(s.ServChain.GetContext(), senderAcc, "aserv")
						Expect(nativeServ.IsZero()).To(BeTrue())
						ibcServ := s.IBCOsmosisChain.GetSimApp().BankKeeper.GetBalance(s.IBCOsmosisChain.GetContext(), receiverAcc, teststypes.AservIbcdenom)
						Expect(ibcServ).To(Equal(sdk.NewCoin(teststypes.AservIbcdenom, coinServ.Amount)))

						// Check that the uosmo were recovered
						ibcOsmo := s.ServChain.App.(*app.Serv).BankKeeper.GetBalance(s.ServChain.GetContext(), receiverAcc, teststypes.UosmoIbcdenom)
						Expect(ibcOsmo.IsZero()).To(BeTrue())
						nativeOsmo := s.IBCOsmosisChain.GetSimApp().BankKeeper.GetBalance(s.IBCOsmosisChain.GetContext(), receiverAcc, "uosmo")
						Expect(nativeOsmo).To(Equal(coinOsmo))

						// Check that the ibc/uatom were retrieved
						osmoIBCAtom := s.ServChain.App.(*app.Serv).BankKeeper.GetBalance(s.ServChain.GetContext(), receiverAcc, teststypes.UatomOsmoIbcdenom)
						Expect(osmoIBCAtom.IsZero()).To(BeTrue())
						ibcAtom := s.IBCOsmosisChain.GetSimApp().BankKeeper.GetBalance(s.IBCOsmosisChain.GetContext(), senderAcc, teststypes.UatomIbcdenom)
						Expect(ibcAtom).To(Equal(sdk.NewCoin(teststypes.UatomIbcdenom, sdk.NewInt(10))))
					})
				})
			})
		})
	})
})
