package lang

import (
	"chatbot/utils/engine_pool"
	"github.com/icyseptember2237/engine"
	"unicode"
)

func init() {
	engine_pool.RegisterModule(func(i engine.Engine) {
		i.RegisterModule(moduleName, moduleMethods)
	})
}

const moduleName = "lang"

var moduleMethods = map[string]interface{}{
	"checkText": checkText,
}

func checkText(text string) string {
	hasEn, hasJa := false, false
	for _, c := range text {
		if unicode.Is(unicode.Hiragana, c) || unicode.Is(unicode.Katakana, c) {
			hasJa = true
		}
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			hasEn = true
		}
		if hasEn && hasJa {
			break
		}
	}
	if hasEn && !hasJa {
		return "zh" // 中英混合
	}
	if !hasEn && hasJa {
		return "all_ja" // 全日文
	}
	if hasEn && hasJa {
		return "auto" // 自动识别
	}
	return "all_zh" // 全中文
}
