package database_test

import (
	feegranttypes "cosmossdk.io/x/feegrant"
	sdk "github.com/cosmos/cosmos-sdk/types"

	dbtypes "github.com/forbole/callisto/v4/database/types"
	"github.com/forbole/callisto/v4/types"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveFeeGrantAllowance() {
	allowance := &feegranttypes.BasicAllowance{SpendLimit: nil, Expiration: nil}
	granter, err := sdk.AccAddressFromBech32("cosmos1ltzt0z992ke6qgmtjxtygwzn36km4cy6cqdknt")
	suite.Require().NoError(err)

	grantee, err := sdk.AccAddressFromBech32("cosmos1re6zjpyczs0w7flrl6uacl0r4teqtyg62crjsn")
	suite.Require().NoError(err)

	feeGrant, err := feegranttypes.NewGrant(granter, grantee, allowance)
	suite.Require().NoError(err)

	// Store the allowance
	err = suite.database.SaveFeeGrantAllowance(types.NewFeeGrant(feeGrant, 121622))
	suite.Require().NoError(err)

	// Test double insertion
	err = suite.database.SaveFeeGrantAllowance(types.NewFeeGrant(feeGrant, 121622))
	suite.Require().NoError(err, "storing existing grant allowance should return no error")

	// Verify the data
	var rows []dbtypes.FeeAllowanceRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM fee_grant_allowance`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().Equal(rows[0].Granter, granter.String())
	suite.Require().Equal(rows[0].Grantee, grantee.String())
	suite.Require().Equal(rows[0].Height, int64(121622))

}

func (suite *DbTestSuite) TestBigDipperDb_RemoveFeeGrantAllowance() {
	allowance := &feegranttypes.BasicAllowance{SpendLimit: nil, Expiration: nil}
	granter, err := sdk.AccAddressFromBech32("cosmos1ltzt0z992ke6qgmtjxtygwzn36km4cy6cqdknt")
	suite.Require().NoError(err)

	grantee, err := sdk.AccAddressFromBech32("cosmos1re6zjpyczs0w7flrl6uacl0r4teqtyg62crjsn")
	suite.Require().NoError(err)

	feeGrant, err := feegranttypes.NewGrant(granter, grantee, allowance)
	suite.Require().NoError(err)

	err = suite.database.SaveFeeGrantAllowance(types.NewFeeGrant(feeGrant, 121622))
	suite.Require().NoError(err)

	// Delete the data
	err = suite.database.DeleteFeeGrantAllowance(types.NewGrantRemoval(
		"cosmos1re6zjpyczs0w7flrl6uacl0r4teqtyg62crjsn",
		"cosmos1ltzt0z992ke6qgmtjxtygwzn36km4cy6cqdknt",
		122222,
	))
	suite.Require().NoError(err)

	// verify the data
	var count int
	err = suite.database.SQL.QueryRow(`SELECT COUNT(*) FROM fee_grant_allowance`).Scan(&count)
	suite.Require().NoError(err)
	suite.Require().Equal(0, count)
}
