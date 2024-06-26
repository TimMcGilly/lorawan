package fragmentation

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFragmentation(t *testing.T) {
	tests := []struct {
		Name                   string
		Command                Command
		Bytes                  []byte
		Uplink                 bool
		ExpectedMarshalError   error
		ExpectedUnmarshalError error
	}{
		{
			Name: "PackageVersionReq",
			Command: Command{
				CID: PackageVersionReq,
			},
			Bytes: []byte{0x00},
		},
		{
			Name:   "PackageVersionAns",
			Uplink: true,
			Command: Command{
				CID: PackageVersionAns,
				Payload: &PackageVersionAnsPayload{
					PackageIdentifier: 1,
					PackageVersion:    1,
				},
			},
			Bytes: []byte{0x00, 0x01, 0x01},
		},
		{
			Name:                   "PackageVersionAns invalid bytes",
			Uplink:                 true,
			Bytes:                  []byte{0x00, 0x01},
			ExpectedUnmarshalError: errors.New("lorawan/applayer/fragmentation: 2 bytes are expected"),
		},
		{
			Name: "FragSessionSetupReq",
			Command: Command{
				CID: FragSessionSetupReq,
				Payload: &FragSessionSetupReqPayload{
					FragSession: FragSessionSetupReqPayloadFragSession{
						FragIndex:      3,
						McGroupBitMask: [4]bool{true, false, true, false},
					},
					NbFrag:   513,
					FragSize: 255,
					Control: FragSessionSetupReqPayloadControl{
						FragmentationMatrix: 5,
						BlockAckDelay:       4,
					},
					Padding:    129,
					Descriptor: [4]byte{0x01, 0x02, 0x03, 0x04},
				},
			},
			Bytes: []byte{0x02, 0x35, 0x01, 0x02, 0xff, 0x2c, 0x81, 0x01, 0x02, 0x03, 0x04},
		},
		{
			Name:                   "FragSessionSetupReq invalid bytes",
			Bytes:                  []byte{0x02, 0x35, 0x01, 0x02, 0xff, 0x2c, 0x81, 0x01, 0x02, 0x03},
			ExpectedUnmarshalError: errors.New("lorawan/applayer/fragmentation: 10 bytes are expected"),
		},
		{
			Name:   "FragSessionSetupAns",
			Uplink: true,
			Command: Command{
				CID: FragSessionSetupAns,
				Payload: &FragSessionSetupAnsPayload{
					StatusBitMask: FragSessionSetupAnsPayloadStatusBitMask{
						FragIndex:                    3,
						WrongDescriptor:              true,
						FragSessionIndexNotSupported: true,
						NotEnoughMemory:              true,
						EncodingUnsupported:          true,
					},
				},
			},
			Bytes: []byte{0x02, 0xcf},
		},
		{
			Name:                   "FragSessionSetupAns invalid bytes",
			Uplink:                 true,
			Bytes:                  []byte{0x02},
			ExpectedUnmarshalError: errors.New("lorawan/applayer/fragmentation: 1 byte is expected"),
		},
		{
			Name: "FragSessionDeleteReq",
			Command: Command{
				CID: FragSessionDeleteReq,
				Payload: &FragSessionDeleteReqPayload{
					Param: FragSessionDeleteReqPayloadParam{
						FragIndex: 3,
					},
				},
			},
			Bytes: []byte{0x03, 0x03},
		},
		{
			Name:                   "FragSessionDeleteReq invalid bytes",
			Bytes:                  []byte{0x03},
			ExpectedUnmarshalError: errors.New("lorawan/applayer/fragmentation: 1 byte is expected"),
		},
		{
			Name:   "FragSessionDeleteAns",
			Uplink: true,
			Command: Command{
				CID: FragSessionDeleteAns,
				Payload: &FragSessionDeleteAnsPayload{
					Status: FragSessionDeleteAnsPayloadStatus{
						FragIndex:           3,
						SessionDoesNotExist: true,
					},
				},
			},
			Bytes: []byte{0x03, 0x07},
		},
		{
			Name:                   "FragSessionDeleteAns invalid bytes",
			Uplink:                 true,
			Bytes:                  []byte{0x03},
			ExpectedUnmarshalError: errors.New("lorawan/applayer/fragmentation: 1 byte is expected"),
		},
		{
			Name: "DataFragment",
			Command: Command{
				CID: DataFragment,
				Payload: &DataFragmentPayload{
					IndexAndN: DataFragmentPayloadIndexAndN{
						FragIndex: 3,
						N:         513,
					},
					Payload: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06},
				},
			},
			Bytes: []byte{0x08, 0x01, 0xc2, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06},
		},
		{
			Name:                   "DataFragment invalid bytes",
			Bytes:                  []byte{0x08, 0x01},
			ExpectedUnmarshalError: errors.New("lorawan/applayer/fragmentation: 2 bytes are expected"),
		},
		{
			Name: "RetransmitDataFragment",
			Command: Command{
				CID: RetransmitDataFragment,
				Payload: &RetransmitDataFragmentPayload{
					IndexAndN: RetransmitDataFragmentPayloadIndexAndN{
						FragIndex: 3,
						N:         513,
					},
					Payload: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06},
				},
			},
			Bytes: []byte{0x05, 0x01, 0xc2, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06},
		},
		{
			Name:                   "RetransmitDataFragment invalid bytes",
			Bytes:                  []byte{0x05, 0x01},
			ExpectedUnmarshalError: errors.New("lorawan/applayer/fragmentation: 2 bytes are expected"),
		},
		{
			Name: "FragSessionStatusReq",
			Command: Command{
				CID: FragSessionStatusReq,
				Payload: &FragSessionStatusReqPayload{
					FragStatusReqParam: FragSessionStatusReqPayloadFragStatusReqParam{
						Participants: true,
						FragIndex:    3,
					},
				},
			},
			Bytes: []byte{0x01, 0x07},
		},
		{
			Name:                   "FragSessionStatusReq invalid bytes",
			Bytes:                  []byte{0x01},
			ExpectedUnmarshalError: errors.New("lorawan/applayer/fragmentation: 1 byte is expected"),
		},
		{
			Name:   "FragSessionStatusAns",
			Uplink: true,
			Command: Command{
				CID: FragSessionStatusAns,
				Payload: &FragSessionStatusAnsPayload{
					ReceivedAndIndex: FragSessionStatusAnsPayloadReceivedAndIndex{
						FragIndex:      3,
						NbFragReceived: 513,
					},
					MissingFrag: 255,
					Status: FragSessionStatusAnsPayloadStatus{
						NotEnoughMatrixMemory: true,
					},
				},
			},
			Bytes: []byte{0x01, 0x01, 0xc2, 0xff, 0x01},
		},
		{
			Name:                   "FragSessionStatusAns invalid bytes",
			Uplink:                 true,
			Bytes:                  []byte{0x01, 0x01, 0xc2, 0xff},
			ExpectedUnmarshalError: errors.New("lorawan/applayer/fragmentation: 4 bytes are expected"),
		},
		{
			Name: "FragSessionMissingReq",
			Command: Command{
				CID: FragSessionMissingReq,
				Payload: &FragSessionMissingReqPayload{
					FragSessionMissingReqParam: FragSessionMissingReqPayloadFragSessionMissingReqParam{
						Participants: true,
						FragIndex:    2,
					},
				},
			},
			Bytes: []byte{0x06, 0x05},
		},
		{
			Name:                   "FragSessionMissingReq invalid bytes",
			Bytes:                  []byte{0x06},
			ExpectedUnmarshalError: errors.New("lorawan/applayer/fragmentation: 1 byte is expected"),
		},
		{
			Name:   "FragSessionMissingBitAns",
			Uplink: true,
			Command: Command{
				CID: FragSessionMissingBitAns,
				Payload: &FragSessionMissingBitAnsPayload{
					MissingAnsHeader: FragSessionMissingBitAnsPayloadMissingAnsHeader{
						FragIndex:          3,
						BitArrayStartIndex: 513,
					},
					NumMissingAns:    3,
					ReceivedBitArray: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06},
				},
			},
			Bytes: []byte{0x07, 0x01, 0xc2, 0x03, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06},
		},
		{
			Name:                   "FragSessionMissingBitAns invalid bytes",
			Uplink:                 true,
			Bytes:                  []byte{0x07, 0x2, 0x2},
			ExpectedUnmarshalError: errors.New("lorawan/applayer/fragmentation: 3 bytes are expected"),
		},
		{
			Name:   "FragSessionMissingListAns",
			Uplink: true,
			Command: Command{
				CID: FragSessionMissingListAns,
				Payload: &FragSessionMissingListAnsPayload{
					MissingAnsHeader: FragSessionMissingListAnsPayloadMissingAnsHeader{
						FragIndex: 3,
					},
					NumMissingAns: 5,
					MissingList:   []uint16{0x0001, 0x0200, 0x3000},
				},
			},
			Bytes: []byte{0x06, 0x3, 0x5, 0x01, 0x00, 0x00, 0x02, 0x00, 0x30},
		},
		{
			Name:                   "FragSessionMissingBitAns invalid bytes",
			Uplink:                 true,
			Bytes:                  []byte{0x06, 0x2},
			ExpectedUnmarshalError: errors.New("lorawan/applayer/fragmentation: 2 bytes are expected"),
		},
	}

	for _, tst := range tests {
		t.Run(tst.Name, func(t *testing.T) {
			assert := require.New(t)

			if tst.ExpectedMarshalError != nil {
				_, err := tst.Command.MarshalBinary()
				assert.Equal(tst.ExpectedMarshalError, err)
			} else if tst.ExpectedUnmarshalError != nil {
				var cmd Command
				err := cmd.UnmarshalBinary(tst.Uplink, tst.Bytes)
				assert.Equal(tst.ExpectedUnmarshalError, err)
			} else {
				cmds := Commands{tst.Command}
				b, err := cmds.MarshalBinary()
				assert.NoError(err)
				assert.Equal(tst.Bytes, b)

				cmds = Commands{}
				assert.NoError(cmds.UnmarshalBinary(tst.Uplink, tst.Bytes))
				assert.Len(cmds, 1)
				assert.Equal(tst.Command, cmds[0])
			}
		})
	}
}
