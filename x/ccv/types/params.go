package types

import (
	fmt "fmt"
	time "time"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

const (
	// about 2 hr at 7.6 seconds per blocks
	DefaultBlocksPerDistributionTransmission = 1000

	// Default transfer timeout period is 1 hour, less than the default blocks
	// per dist transmission * average block time.
	// Since IBC token transfers do not have to be in order, it could be easier
	// to reason about the distribution protocol if the previous reward times out
	// before sending the next one. Note that on timeout, the transferred funds are
	// added back to the pool, so the next transfer will include them as well.
	DefaultTransferTimeoutPeriod = time.Hour

	// The default fraction of tokens allocated to the consumer redistribution address
	// during distribution events. The fraction is a string representing a
	// decimal number. For example "0.75" would represent 75%.
	DefaultConsumerRedistributeFrac = "0.75"

	// Default number of historical info entries to persist in store.
	// We use the same default as the staking module, but use a signed integer
	// so that negative values can be caught during parameter validation in a readable way,
	// (and for consistency with other protobuf schemas defined for ccv).
	DefaultHistoricalEntries = int64(stakingtypes.DefaultHistoricalEntries)

	// In general, the default unbonding period on the consumer is one day less
	// than the default unbonding period on the provider, where the provider uses
	// the staking module default.
	DefaultConsumerUnbondingPeriod = stakingtypes.DefaultUnbondingTime - 24*time.Hour

	// By default, the bottom 5% of the validator set can opt out of validating consumer chains
	DefaultSoftOptOutThreshold = "0.05"
)

// Reflection based keys for params subspace
var (
	KeyEnabled                           = []byte("Enabled")
	KeyBlocksPerDistributionTransmission = []byte("BlocksPerDistributionTransmission")
	KeyDistributionTransmissionChannel   = []byte("DistributionTransmissionChannel")
	KeyProviderFeePoolAddrStr            = []byte("ProviderFeePoolAddrStr")
	KeyTransferTimeoutPeriod             = []byte("TransferTimeoutPeriod")
	KeyConsumerRedistributionFrac        = []byte("ConsumerRedistributionFraction")
	KeyHistoricalEntries                 = []byte("HistoricalEntries")
	KeyConsumerUnbondingPeriod           = []byte("UnbondingPeriod")
	KeySoftOptOutThreshold               = []byte("SoftOptOutThreshold")
	KeyRewardDenoms                      = []byte("RewardDenoms")
	KeyProviderRewardDenoms              = []byte("ProviderRewardDenoms")
)

// ParamKeyTable type declaration for parameters
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates new consumer parameters with provided arguments
func NewParams(enabled bool, blocksPerDistributionTransmission int64,
	distributionTransmissionChannel, providerFeePoolAddrStr string,
	ccvTimeoutPeriod, transferTimeoutPeriod time.Duration,
	consumerRedistributionFraction string, historicalEntries int64,
	consumerUnbondingPeriod time.Duration, softOptOutThreshold string, rewardDenoms, providerRewardDenoms []string,
) Params {
	return Params{
		Enabled:                           enabled,
		BlocksPerDistributionTransmission: blocksPerDistributionTransmission,
		DistributionTransmissionChannel:   distributionTransmissionChannel,
		ProviderFeePoolAddrStr:            providerFeePoolAddrStr,
		CcvTimeoutPeriod:                  ccvTimeoutPeriod,
		TransferTimeoutPeriod:             transferTimeoutPeriod,
		ConsumerRedistributionFraction:    consumerRedistributionFraction,
		HistoricalEntries:                 historicalEntries,
		UnbondingPeriod:                   consumerUnbondingPeriod,
		SoftOptOutThreshold:               softOptOutThreshold,
		RewardDenoms:                      rewardDenoms,
		ProviderRewardDenoms:              providerRewardDenoms,
	}
}

// DefaultParams is the default params for the consumer module
func DefaultParams() Params {
	var rewardDenoms []string
	var provideRewardDenoms []string
	return NewParams(
		false,
		DefaultBlocksPerDistributionTransmission,
		"",
		"",
		DefaultCCVTimeoutPeriod,
		DefaultTransferTimeoutPeriod,
		DefaultConsumerRedistributeFrac,
		DefaultHistoricalEntries,
		DefaultConsumerUnbondingPeriod,
		DefaultSoftOptOutThreshold,
		rewardDenoms,
		provideRewardDenoms,
	)
}

// Validate all ccv-consumer module parameters
func (p Params) Validate() error {
	if err := ValidateBool(p.Enabled); err != nil {
		return err
	}
	if err := ValidatePositiveInt64(p.BlocksPerDistributionTransmission); err != nil {
		return err
	}
	if err := ValidateDistributionTransmissionChannel(p.DistributionTransmissionChannel); err != nil {
		return err
	}
	if err := ValidateProviderFeePoolAddrStr(p.ProviderFeePoolAddrStr); err != nil {
		return err
	}
	if err := ValidateDuration(p.CcvTimeoutPeriod); err != nil {
		return err
	}
	if err := ValidateDuration(p.TransferTimeoutPeriod); err != nil {
		return err
	}
	if err := ValidateStringFraction(p.ConsumerRedistributionFraction); err != nil {
		return err
	}
	if err := ValidatePositiveInt64(p.HistoricalEntries); err != nil {
		return err
	}
	if err := ValidateDuration(p.UnbondingPeriod); err != nil {
		return err
	}
	if err := ValidateSoftOptOutThreshold(p.SoftOptOutThreshold); err != nil {
		return err
	}
	if err := ValidateDenoms(p.RewardDenoms); err != nil {
		return err
	}
	if err := ValidateDenoms(p.ProviderRewardDenoms); err != nil {
		return err
	}
	return nil
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyEnabled, p.Enabled, ValidateBool),
		paramtypes.NewParamSetPair(KeyBlocksPerDistributionTransmission,
			p.BlocksPerDistributionTransmission, ValidatePositiveInt64),
		paramtypes.NewParamSetPair(KeyDistributionTransmissionChannel,
			p.DistributionTransmissionChannel, ValidateDistributionTransmissionChannel),
		paramtypes.NewParamSetPair(KeyProviderFeePoolAddrStr,
			p.ProviderFeePoolAddrStr, ValidateProviderFeePoolAddrStr),
		paramtypes.NewParamSetPair(KeyCCVTimeoutPeriod,
			p.CcvTimeoutPeriod, ValidateDuration),
		paramtypes.NewParamSetPair(KeyTransferTimeoutPeriod,
			p.TransferTimeoutPeriod, ValidateDuration),
		paramtypes.NewParamSetPair(KeyConsumerRedistributionFrac,
			p.ConsumerRedistributionFraction, ValidateStringFraction),
		paramtypes.NewParamSetPair(KeyHistoricalEntries,
			p.HistoricalEntries, ValidatePositiveInt64),
		paramtypes.NewParamSetPair(KeyConsumerUnbondingPeriod,
			p.UnbondingPeriod, ValidateDuration),
		paramtypes.NewParamSetPair(KeySoftOptOutThreshold,
			p.SoftOptOutThreshold, ValidateSoftOptOutThreshold),
		paramtypes.NewParamSetPair(KeyRewardDenoms,
			p.RewardDenoms, ValidateDenoms),
		paramtypes.NewParamSetPair(KeyProviderRewardDenoms,
			p.ProviderRewardDenoms, ValidateDenoms),
	}
}

func ValidateProviderFeePoolAddrStr(i interface{}) error {
	// Accept empty string as valid, since this will be the default value on genesis
	if i == "" {
		return nil
	}
	// Otherwise validate as usual for a bech32 address
	return ValidateBech32(i)
}

func ValidateSoftOptOutThreshold(i interface{}) error {
	str, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	dec, err := sdktypes.NewDecFromStr(str)
	if err != nil {
		return err
	}
	if dec.IsNegative() {
		return fmt.Errorf("soft opt out threshold cannot be negative, got %s", str)
	}
	if !dec.Sub(sdktypes.MustNewDecFromStr("0.2")).IsNegative() {
		return fmt.Errorf("soft opt out threshold cannot be greater than 0.2, got %s", str)
	}
	return nil
}

func ValidateDenoms(i interface{}) error {
	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// iterate over the denoms, turning them into coins and validating them
	for _, denom := range v {
		coin := sdktypes.Coin{
			Denom:  denom,
			Amount: sdktypes.NewInt(0),
		}

		if err := coin.Validate(); err != nil {
			return err
		}
	}

	return nil
}
