package types_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"praios.lf-net.org/littlefox/gotoshock/pkg/types"
)

func messageFromString(data string) types.Message {
	ret := types.Message{}

	Expect(data).To(HaveLen(len(ret)))

	for i, c := range data {
		ret[i] = c == '1'
	}
	return ret
}

var _ = Describe("Message", func() {
	channelFormat := func(a, b string) string {
		return fmt.Sprintf("00%s000000000000000000000000000000%s00", a, b)
	}

	operationFormat := func(a, b string) string {
		return fmt.Sprintf("000000%s000000000000000000000000%s000000", a, b)
	}

	intensityFormat := func(intensity string) string {
		return fmt.Sprintf("00000000000000000000000000%s000000000", intensity)
	}

	Context("Build", func() {
		It("adds the constants as expected", func() {
			expected := messageFromString("010000000001011100010101100000000000000000")
			msg := new(types.Message).
				Build()

			Expect(*msg).To(Equal(expected))
		})
	})

	DescribeTable("SetChannel",
		func(ch types.Channel, expectedMessage string) {
			expected := messageFromString(expectedMessage)

			msg := new(types.Message).
				SetChannel(ch)

			Expect(*msg).To(Equal(expected))
		},
		Entry("1", types.Channel1, channelFormat("0000", "1111")),
		Entry("2", types.Channel2, channelFormat("1110", "1000")),
		Entry("Unknown", types.Channel(12), channelFormat("1100", "1100")),
	)

	DescribeTable("SetOperation",
		func(op types.Operation, expectedMessage string) {
			expected := messageFromString(expectedMessage)

			msg := new(types.Message).
				SetOperation(op)

			Expect(*msg).To(Equal(expected))
		},
		Entry("Shock", types.OperationShock, operationFormat("001", "011")),
		Entry("Vibrate", types.OperationVibrate, operationFormat("010", "101")),
		Entry("Beep", types.OperationBeep, operationFormat("100", "110")),
		Entry("Unknown", types.Operation(7), operationFormat("111", "000")),
	)

	DescribeTable("SetIntensity",
		func(intensity int, expectedMessage string) {
			expected := messageFromString(expectedMessage)

			msg := new(types.Message).
				SetIntensity(types.Intensity(intensity))

			Expect(*msg).To(Equal(expected))
		},
		Entry("38", 38, intensityFormat("0100110")),
		Entry("10", 10, intensityFormat("0001010")),
		Entry("0", 0, intensityFormat("0000000")),
		Entry("50", 50, intensityFormat("0110010")),
	)

	DescribeTable("Full",
		func(ch types.Channel, op types.Operation, intensity int, expectedMessage string) {
			expected := messageFromString(expectedMessage)

			msg := types.NewMessage().
				SetChannel(ch).
				SetOperation(op).
				SetIntensity(types.Intensity(intensity)).
				Build()

			Expect(*msg).To(Equal(expected))
		},
		Entry("Channel 1 Shock 38", types.Channel1, types.OperationShock, 38, "010000001001011100010101100100110011111100"),
		Entry("Channel 2 Shock 10", types.Channel2, types.OperationShock, 10, "011110001001011100010101100001010011100000"),
		Entry("Channel 1 Shock 0", types.Channel1, types.OperationShock, 0, "010000001001011100010101100000000011111100"),
		Entry("Channel 1 Vibrate 50", types.Channel1, types.OperationVibrate, 50, "010000010001011100010101100110010101111100"),
		Entry("Channel 2 Beep 0", types.Channel2, types.OperationBeep, 0, "011110100001011100010101100000000110100000"),
	)
})
