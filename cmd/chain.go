package cmd

import (
	"fmt"
	"log"

	"github.com/jiangjincc/islands/block"

	"github.com/spf13/cobra"
)

var (
	data     string
	addresss string

	chainCmd = &cobra.Command{
		Use:   "chain",
		Short: "区块链命令行工具",
		Long:  ``,
	}

	chainAddCmd = &cobra.Command{
		Use:   "add",
		Short: "添加区块",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			b := block.GetBlockchain()
			err := b.AddBlockToBlockChain([]*block.Transaction{})
			if err != nil {
				log.Panic(err)
			}
		},
	}

	chainListCmd = &cobra.Command{
		Use:   "list",
		Short: "列出所有区块",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			block := block.GetBlockchain()
			fmt.Println("打印区块")
			block.PrintBlocks()
		},
	}

	chainInitCmd = &cobra.Command{
		Use:   "init",
		Short: "初始化区块链",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			block.CreateBlockchainWithGenesisBlock(addresss)

			blockchain := block.GetBlockchain()
			utxoSet := &block.UTXOSet{Blockchain: blockchain}
			utxoSet.ResetUTXOSet()
		},
	}
	chainTestCmd = &cobra.Command{
		Use:   "test",
		Short: "测试",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			//block := block.GetBlockchain()
			//utxoMap := block.FindUTXOMap()
			//fmt.Println(utxoMap)
			blockchain := block.GetBlockchain()
			utxoSet := &block.UTXOSet{Blockchain: blockchain}
			utxoSet.ResetUTXOSet()
		},
	}
)

func chainCmdExecute(rootCmd *cobra.Command) {
	// 解析参数
	chainAddCmd.Flags().StringVarP(&data, "data", "d", "", "区块数据 (required)")
	_ = chainAddCmd.MarkFlagRequired("data")

	chainInitCmd.Flags().StringVarP(&addresss, "address", "a", "", "地址(required)")
	_ = chainAddCmd.MarkFlagRequired("address")

	rootCmd.AddCommand(chainCmd)
	chainCmd.AddCommand(chainAddCmd)
	chainCmd.AddCommand(chainListCmd)
	chainCmd.AddCommand(chainInitCmd)
	chainCmd.AddCommand(chainTestCmd)
}
