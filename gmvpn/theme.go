/*
 * @Author: liuwei lyy9645@163.com
 * @Date: 2023-05-19 23:34:35
 * @LastEditors: liuwei lyy9645@163.com
 * @LastEditTime: 2023-05-20 22:06:50
 * @FilePath: /gmvpn/ui/theme.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type myTheme struct{}

var _ fyne.Theme = (*myTheme)(nil)

// return bundled font resource
func (*myTheme) Font(s fyne.TextStyle) fyne.Resource {
	if s.Monospace {
		return theme.DefaultTheme().Font(s)
	}
	if s.Bold {
		if s.Italic {
			return theme.DefaultTheme().Font(s)
		}
		return resourceAlibabaPuHuiTi255RegularTtf
	}
	if s.Italic {
		return theme.DefaultTheme().Font(s)
	}
	return resourceAlibabaPuHuiTi255RegularTtf

}

func (*myTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, v)
}

func (*myTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (*myTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}
