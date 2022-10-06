package encryptutil_test

import (
	"fmt"
	"testing"

	"github.com/rabee-inc/go-pkg/encryptutil"
	"github.com/rabee-inc/go-pkg/randutil"
)

func TestEventApplyStrategyEveryone(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		key := "+b=$3JDmDseH'w}LM$7EMWkqy3_&d=kk"
		texts := []string{
			"",
			"a",
			"あ",
			"ア",
			"ｱｲｳ",
			"123",
			"<>?+*`!#$%&'()\n\t\"",
			"＜＞？；：」「＠！”＃＄％＆’（）＿・",
			"0ROlNc31oEFMunR5f4Dm4NPuTFEDt5vuumv2ScWSWC5MDQAfU4Xrik2DiN4dTgKtVAaf8nulw2EhgrkQnu6cdQaCjvoCvXWow10TMjPO0Mp0Gd85J",
			"蓬滷棉杗苠浢萎樼淦棌ブ灆萪渶Ｏ潳薕U6橅溡欛茽檤ｲ浬楨梼椒を洬薀蒲椓渥洜ぼ蒁o枊蘗棉槺潍ヵ椑蘚ムﾁ枱藲荥+渊葃苵ク薊湅櫑栜潵泒汊氺槠(汔椙荅汕潣瀯蒄泖深果溧榌枮G栁滛澭荫著涓濶蘋溧欙槷潾漂莗桛楳1芵泲",
			"あab🈵⛄あいうｱｲｳ123１２３",
			"🌱 🌲 🌳 🌴 🌵 🌷 🌸 🌹 🌺 🌻 🌼 🌽 🌾 🌿 🍀 🍁 🍂 🍃 🍄 🍅 🍆 🍇 🍈 🍉 🍊 🍋 🍌 🍍 🍎 🍏 🍐 🍑 🍒 🍓",
			"'FmgMQpB:+E;6MBr?%8Z!?T*Recy,_ME#S;dnka>g4d]Mr|4hlO(u(^7K~M,cqFN#-0E<KXR>bBQ^*~T)DR~E&J*w;m?j`L*TQ7G",
			"東京都◯◯区✗✗ １−２−３ ほげほげハイツ２００号室",
		}
		for i := 0; i < 1000; i++ {
			text, err := randutil.String(randutil.Int(1, 200))
			if err != nil {
				t.Error(err)
				return
			}
			texts = append(texts, text)
		}

		for _, text := range texts {
			encText, err := encryptutil.Encrypt(text, key)
			if err != nil {
				t.Error(err)
				return
			}
			fmt.Println(encText)
			decText, err := encryptutil.Decrypt(encText, key)
			if text != decText {
				t.Errorf("no match text: %s != %s", text, decText)
				return
			}
		}
	})
}
