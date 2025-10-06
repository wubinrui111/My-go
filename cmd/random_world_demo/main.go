package main

import (
	"fmt"
	"math/rand"
	"time"

	"mygo/internal/pkg/game"
)

func main() {
	fmt.Println("演示随机世界生成:")

	// 设置随机种子
	rand.Seed(time.Now().UnixNano())

	// 生成几个不同的世界并比较它们的方块数量
	for i := 1; i <= 3; i++ {
		fmt.Printf("\n第 %d 次生成世界:\n", i)

		// 创建新游戏实例（会自动生成随机世界）
		g := game.NewGame()
		_ = g
		fmt.Println("成功创建了一个随机世界!")
		fmt.Println("每次运行程序都会生成不同的地形和环境!")
	}

	fmt.Println("\n现在可以运行主游戏来体验随机生成的世界了!")
}
