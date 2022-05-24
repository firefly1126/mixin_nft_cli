package cmd

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"
	"time"

	"mixin_nft_cli/trident"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/qrcode"
	"github.com/fox-one/pkg/uuid"

	"github.com/fox-one/mixin-sdk-go/nft"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var assetCmd = cobra.Command{
	Use: "asset",
	RunE: func(cmd *cobra.Command, args []string) error {
		rsp, err := mixin.GetRestyClient().NewRequest().Get("/network/assets/c94ac88f-4671-3976-b60a-09064f1811e8")
		if err != nil {
			return err
		}

		return prettyPrint(rsp.Body())
	},
}

var outputCmd = cobra.Command{
	Use: "output",
	RunE: func(cmd *cobra.Command, args []string) error {
		mixin.GenerateCollectibleTokenID("", 1)
		ctx := cmd.Context()
		outputs, err := mixin.NewFromAccessToken(cfg.User.Token).ReadCollectibleOutputs(ctx, []string{cfg.User.UserID}, 1, "", time.Time{}, 500)
		if err != nil {
			return err
		}

		return prettyPrint(outputs)
	},
}

var tokenCmd = cobra.Command{
	Use: "token",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		tokenID, err := cmd.Flags().GetString("token")
		if err != nil {
			return err
		}

		token, err := cfg.Mixin.Client().ReadCollectiblesToken(ctx, tokenID) //mixin.NewFromAccessToken(cfg.User.Token).ReadCollectiblesToken(ctx, tokenID)
		if err != nil {
			return err
		}

		return prettyPrint(token)
	},
}

var metaDataCmd = cobra.Command{
	Use: "meta",
	RunE: func(cmd *cobra.Command, args []string) error {
		metaHash, err := cmd.Flags().GetString("hash")
		if err != nil {
			return err
		}

		metaData, err := trident.GetMetaData(metaHash)
		if err != nil {
			return err
		}

		return prettyPrint(metaData)
	},
}

var collectionCmd = cobra.Command{
	Use: "collection",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		collectionID, err := cmd.Flags().GetString("id")
		if err != nil {
			return err
		}

		collection, err := mixin.NewFromAccessToken(cfg.User.Token).ReadCollectibleCollection(ctx, collectionID) //mixin.GetRestyClient().NewRequest().SetAuthToken(cfg.User.Token).Get("/collectibles/collections/" + collectionID)
		if err != nil {
			return err
		}

		return prettyPrint(collection)
	},
}

var transferNFTCmd = cobra.Command{
	Use: "transfer",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		tokenID, err := cmd.Flags().GetString("token")
		if err != nil {
			return err
		}

		receiver, err := cmd.Flags().GetString("receiver")
		if err != nil {
			return err
		}

		// get unspent outputs
		var outputTarget *mixin.CollectibleOutput
		os, err := mixin.NewFromAccessToken(cfg.User.Token).ReadCollectibleOutputs(ctx, []string{cfg.User.UserID}, 1, "", time.Time{}, 500)
		if err != nil {
			return err
		}
		for _, o := range os {
			if o.TokenID == tokenID && o.State == "unspent" {
				outputTarget = o
			}
		}
		if outputTarget == nil {
			return errors.New("unspent output not found")
		}

		// get token
		token, err := mixin.NewFromAccessToken(cfg.User.Token).ReadCollectiblesToken(ctx, tokenID)
		if err != nil {
			return err
		}

		// make collectible transaction
		tran, err := mixin.NewFromAccessToken(cfg.User.Token).MakeCollectibleTransaction(ctx, outputTarget, token, []string{receiver}, 1)
		if err != nil {
			return err
		}
		signedTx, err := tran.DumpTransaction()
		if err != nil {
			return err
		}

		// create collectible request
		collectibleRequest, err := mixin.NewFromAccessToken(cfg.User.Token).CreateCollectibleRequest(ctx, mixin.CollectibleRequestActionSign, signedTx)
		if err != nil {
			return err
		}

		// sign collectible request
		url := mixin.URL.Codes(collectibleRequest.CodeID)
		logrus.Infoln(url)
		qrcode.Fprint(cmd.OutOrStdout(), url)

		return nil
	},
}

var spendNFTCmd = cobra.Command{
	Use: "spend",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		tokenID, err := cmd.Flags().GetString("token")
		if err != nil {
			return err
		}

		// get unspent outputs
		var outputTarget *mixin.CollectibleOutput
		os, err := mixin.NewFromAccessToken(cfg.User.Token).ReadCollectibleOutputs(ctx, []string{cfg.User.UserID}, 1, "", time.Time{}, 500)
		if err != nil {
			return err
		}
		for _, o := range os {
			if o.TokenID == tokenID && o.State == "signed" {
				outputTarget = o
			}
		}
		if outputTarget == nil {
			return errors.New("unspent output not found")
		}

		tx, err := mixin.TransactionFromRaw(outputTarget.SignedTx)
		if err != nil {
			return err
		}

		logrus.Infoln("aggregateSignature:", tx.AggregatedSignature)
		if tx.AggregatedSignature == nil {
			return errors.New("no aggregatedSignature")
		}

		hs, err := mixin.NewFromAccessToken(cfg.User.Token).SendRawTransaction(ctx, outputTarget.SignedTx)
		if err != nil {
			return err
		}

		logrus.Infoln(hs)

		return nil
	},
}

var mintNFTCmd = cobra.Command{
	Use: "mint",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		metaDataPath, err := cmd.Flags().GetString("metadata")
		if err != nil {
			return err
		}
		hash, err := cmd.Flags().GetString("hash")
		if err != nil {
			return err
		}

		//upload data to trident
		bts, err := ioutil.ReadFile(metaDataPath)
		if err != nil {
			return err
		}

		var metaData trident.MetaData
		if err := json.Unmarshal(bts, &metaData); err != nil {
			return err
		}

		payload := trident.CreateMetaDataPayload{
			MetaData: metaData,
			MetaHash: hash,
		}

		rsp, err := trident.CreateMetaData(cfg.User.Token, &payload)
		if err != nil {
			return err
		}

		prettyPrint(rsp)

		//push data to mixin network
		hs, err := mixin.HashFromString(hash)
		if err != nil {
			return err
		}

		tokenID, err := strconv.ParseInt(metaData.Token.ID, 10, 64)
		if err != nil {
			return err
		}

		nfo := nft.BuildMintNFO(metaData.Collection.ID, tokenID, hs)
		memo := base64.RawURLEncoding.EncodeToString(nfo)
		traceID := uuid.New()
		input := mixin.TransferInput{
			AssetID: nft.MintAssetId,
			Amount:  nft.MintMinimumCost,
			TraceID: traceID,
			Memo:    memo,
		}

		input.OpponentMultisig.Receivers = nft.GroupMembers
		input.OpponentMultisig.Threshold = nft.GroupThreshold
		payment, err := mixin.NewFromAccessToken(cfg.User.Token).VerifyPayment(ctx, input)
		if err != nil {
			return err
		}

		url := mixin.URL.Codes(payment.CodeID)
		logrus.Infoln(url)
		qrcode.Fprint(cmd.OutOrStdout(), url)

		return nil
	},
}

var metaHashCmd = cobra.Command{
	Use:  "metahash",
	Long: "metahash --fields aaaa,bbbb,ccccc",
	RunE: func(cmd *cobra.Command, args []string) error {
		fields, err := cmd.Flags().GetStringSlice("fields")
		if err != nil {
			return err
		}

		content := ""
		for _, v := range fields {
			content += v
		}

		logrus.Infoln("content:", content)

		hh := mixin.NewHash([]byte(content))
		logrus.Infoln("hash:", hh.String())

		return nil
	},
}

func init() {
	root.AddCommand(&assetCmd)
	root.AddCommand(&outputCmd)

	tokenCmd.Flags().StringP("token", "t", "", "nft token id")
	root.AddCommand(&tokenCmd)

	metaDataCmd.Flags().StringP("hash", "", "", "meta hash")
	root.AddCommand(&metaDataCmd)

	collectionCmd.Flags().StringP("id", "i", "", "collection id")
	root.AddCommand(&collectionCmd)

	transferNFTCmd.Flags().StringP("token", "t", "", "nft token id")
	transferNFTCmd.Flags().StringP("receiver", "r", "", "receiver user id")
	root.AddCommand(&transferNFTCmd)

	spendNFTCmd.Flags().StringP("token", "t", "", "nft token id")
	root.AddCommand(&spendNFTCmd)

	mintNFTCmd.Flags().StringP("metadata", "m", "", "metadata file path, e.g.,metadata.example.json")
	mintNFTCmd.Flags().StringP("hash", "h", "", "meta data hash")
	root.AddCommand(&mintNFTCmd)

	metaHashCmd.Flags().StringSliceP("fields", "f", []string{}, "meta hash fields, e.g., aaa,bbb,ccc")
	root.AddCommand(&metaHashCmd)
}
