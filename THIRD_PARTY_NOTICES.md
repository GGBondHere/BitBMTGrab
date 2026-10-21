# 第三方版权与许可证声明

发行程序所含的第三方组件及其许可如下。再分发时请保留相应版权与许可声明。

组件版本对应 `go.mod`，Go 标准库许可对应发行使用的 Go 1.26.3。

## 原始项目

本项目基于 [YaphetLee2002/BitBMTGrab](https://github.com/YaphetLee2002/BitBMTGrab) 的初始实现。原仓库未声明许可证，其原始代码不适用本仓库为新增与修改部分提供的 MIT 许可；原始代码的再许可需另行取得作者授权。

## MIT 许可组件

下列组件分别保留表中的版权声明，共同适用本节末尾的完整 MIT 许可文本。组件名称链接到对应版本的上游源码。

| 组件 | 版本 | 版权或作者声明 |
| --- | --- | --- |
| [github.com/aymanbagabas/go-osc52/v2](https://github.com/aymanbagabas/go-osc52/tree/v2.0.1) | `v2.0.1` | Copyright (c) 2022 Ayman Bagabas |
| [github.com/catppuccin/go](https://github.com/catppuccin/go/tree/v0.3.0) | `v0.3.0` | Copyright (c) 2021 Catppuccin |
| [github.com/charmbracelet/bubbles](https://github.com/charmbracelet/bubbles/tree/23b8fd6302d7) | `v0.21.1-0.20250623103423-23b8fd6302d7` | Copyright (c) 2020-2023 Charmbracelet, Inc |
| [github.com/charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea/tree/v1.3.6) | `v1.3.6` | Copyright (c) 2020-2023 Charmbracelet, Inc |
| [github.com/charmbracelet/colorprofile](https://github.com/charmbracelet/colorprofile/tree/f60798e515dc) | `v0.2.3-0.20250311203215-f60798e515dc` | Copyright (c) 2020-2024 Charmbracelet, Inc |
| [github.com/charmbracelet/huh](https://github.com/charmbracelet/huh/tree/v0.8.0) | `v0.8.0` | Copyright (c) 2023 Charm |
| [github.com/charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss/tree/v1.1.0) | `v1.1.0` | Copyright (c) 2021-2023 Charmbracelet, Inc |
| [github.com/charmbracelet/x/ansi](https://github.com/charmbracelet/x/tree/ansi/v0.9.3/ansi) | `v0.9.3` | Copyright (c) 2023 Charmbracelet, Inc. |
| [github.com/charmbracelet/x/cellbuf](https://github.com/charmbracelet/x/tree/cellbuf/v0.0.13/cellbuf) | `v0.0.13` | Copyright (c) 2023 Charmbracelet, Inc. |
| [github.com/charmbracelet/x/exp/strings](https://github.com/charmbracelet/x/tree/212f7b056ed0/exp/strings) | `v0.0.0-20240722160745-212f7b056ed0` | Copyright (c) 2023 Charmbracelet, Inc. |
| [github.com/charmbracelet/x/term](https://github.com/charmbracelet/x/tree/term/v0.2.1/term) | `v0.2.1` | Copyright (c) 2023 Charmbracelet, Inc. |
| [github.com/dustin/go-humanize](https://github.com/dustin/go-humanize/tree/v1.0.1) | `v1.0.1` | Copyright (c) 2005-2008 Dustin Sallings &lt;dustin@spy.net&gt; |
| [github.com/erikgeiser/coninput](https://github.com/erikgeiser/coninput/tree/1c3628e74d0f) | `v0.0.0-20211004153227-1c3628e74d0f` | Copyright (c) 2021 Erik G. |
| [github.com/lucasb-eyer/go-colorful](https://github.com/lucasb-eyer/go-colorful/tree/v1.2.0) | `v1.2.0` | Copyright (c) 2013 Lucas Beyer |
| [github.com/mattn/go-isatty](https://github.com/mattn/go-isatty/tree/v0.0.20) | `v0.0.20` | Copyright (c) Yasuhiro MATSUMOTO &lt;mattn.jp@gmail.com&gt; |
| [github.com/mattn/go-localereader](https://github.com/mattn/go-localereader/tree/v0.0.1) | `v0.0.1` | Yasuhiro Matsumoto (a.k.a. mattn)；上游 README 声明使用 MIT |
| [github.com/mattn/go-runewidth](https://github.com/mattn/go-runewidth/tree/v0.0.16) | `v0.0.16` | Copyright (c) 2016 Yasuhiro Matsumoto |
| [github.com/mitchellh/hashstructure/v2](https://github.com/mitchellh/hashstructure/tree/v2.0.2) | `v2.0.2` | Copyright (c) 2016 Mitchell Hashimoto |
| [github.com/muesli/ansi](https://github.com/muesli/ansi/tree/276c6243b2f6) | `v0.0.0-20230316100256-276c6243b2f6` | Copyright (c) 2021 Christian Muehlhaeuser |
| [github.com/muesli/cancelreader](https://github.com/muesli/cancelreader/tree/v0.2.2) | `v0.2.2` | Copyright (c) 2022 Erik Geiser and Christian Muehlhaeuser |
| [github.com/muesli/termenv](https://github.com/muesli/termenv/tree/v0.16.0) | `v0.16.0` | Copyright (c) 2019 Christian Muehlhaeuser |
| [github.com/rivo/uniseg](https://github.com/rivo/uniseg/tree/v0.4.7) | `v0.4.7` | Copyright (c) 2019 Oliver Kuederle |
| [github.com/xo/terminfo](https://github.com/xo/terminfo/tree/abceb7e1c41e) | `v0.0.0-20220910002029-abceb7e1c41e` | Copyright (c) 2016 Anmol Sethi |

Go 的 amd64 runtime 还包含 [Inferno 衍生的 memmove 实现](https://github.com/golang/go/blob/go1.26.3/src/runtime/memmove_amd64.s)，以下版权声明同样适用本节的 MIT 许可：

```text
Copyright © 1994-1999 Lucent Technologies Inc. All rights reserved.
Revisions Copyright © 2000-2007 Vita Nuova Holdings Limited (www.vitanuova.com). All rights reserved.
Portions Copyright 2009 The Go Authors. All rights reserved.
```

```text
MIT License

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

## Go runtime、标准库及 golang.org/x 组件

| 组件 | 版本 |
| --- | --- |
| [Go runtime / 标准库](https://github.com/golang/go/tree/go1.26.3) | `go1.26.3` |
| [golang.org/x/sync](https://github.com/golang/sync/tree/v0.15.0) | `v0.15.0` |
| [golang.org/x/sys](https://github.com/golang/sys/tree/v0.33.0) | `v0.33.0` |
| [golang.org/x/text](https://github.com/golang/text/tree/v0.23.0) | `v0.23.0` |

上述组件适用以下 BSD 许可；Go 中另有声明的部分适用其对应条款。

```text
Copyright 2009 The Go Authors.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google LLC nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
```

相关专利授权：

```text
Additional IP Rights Grant (Patents)

"This implementation" means the copyrightable works distributed by
Google as part of the Go project.

Google hereby grants to You a perpetual, worldwide, non-exclusive,
no-charge, royalty-free, irrevocable (except as stated in this section)
patent license to make, have made, use, offer to sell, sell, import,
transfer and otherwise run, modify and propagate the contents of this
implementation of Go, where such license applies only to those patent
claims, both currently owned or controlled by Google and acquired in
the future, licensable by Google that are necessarily infringed by this
implementation of Go. This grant does not include claims that would be
infringed only as a consequence of further modification of this
implementation. If you or your agent or exclusive licensee institute or
order or agree to the institution of patent litigation against any
entity (including a cross-claim or counterclaim in a lawsuit) alleging
that this implementation of Go or any code incorporated within this
implementation of Go constitutes direct or contributory patent
infringement, or inducement of patent infringement, then any patent
rights granted to you under this License for this implementation of Go
shall terminate as of the date such litigation is filed.
```

## github.com/atotto/clipboard

版本：`v0.1.4` · [上游源码](https://github.com/atotto/clipboard/tree/v0.1.4)

```text
Copyright (c) 2013 Ato Araki. All rights reserved.
Copyright 2013 @atotto. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of @atotto. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
```

## Unicode 数据

`github.com/rivo/uniseg` 的字符属性表使用 Unicode 15.0.0 数据，保留其数据许可；该库自身的 Go 代码仍适用前述 MIT 许可。

来源：[Unicode 15.0.0](https://www.unicode.org/Public/15.0.0/ucd/)、[uniseg 数据来源说明](https://github.com/rivo/uniseg/blob/v0.4.7/graphemeproperties.go)、[同期 Unicode 许可文本](https://github.com/unicode-org/icu/blob/release-73-2/icu4c/LICENSE)。以下仅收录 Unicode 数据许可，不表示本项目包含 ICU。

```text
UNICODE, INC. LICENSE AGREEMENT - DATA FILES AND SOFTWARE

See Terms of Use <https://www.unicode.org/copyright.html>
for definitions of Unicode Inc.'s Data Files and Software.
NOTICE TO USER: Carefully read the following legal agreement.
BY DOWNLOADING, INSTALLING, COPYING OR OTHERWISE USING UNICODE INC.'S
DATA FILES ("DATA FILES"), AND/OR SOFTWARE ("SOFTWARE"),
YOU UNEQUIVOCALLY ACCEPT, AND AGREE TO BE BOUND BY, ALL OF THE
TERMS AND CONDITIONS OF THIS AGREEMENT.
IF YOU DO NOT AGREE, DO NOT DOWNLOAD, INSTALL, COPY, DISTRIBUTE OR USE
THE DATA FILES OR SOFTWARE.

COPYRIGHT AND PERMISSION NOTICE
Copyright © 1991-2023 Unicode, Inc. All rights reserved.
Distributed under the Terms of Use in https://www.unicode.org/copyright.html.
Permission is hereby granted, free of charge, to any person obtaining
a copy of the Unicode data files and any associated documentation
(the "Data Files") or Unicode software and any associated documentation
(the "Software") to deal in the Data Files or Software
without restriction, including without limitation the rights to use,
copy, modify, merge, publish, distribute, and/or sell copies of
the Data Files or Software, and to permit persons to whom the Data Files
or Software are furnished to do so, provided that either
(a) this copyright and permission notice appear with all copies
of the Data Files or Software, or
(b) this copyright and permission notice appear in associated
Documentation.
THE DATA FILES AND SOFTWARE ARE PROVIDED "AS IS", WITHOUT WARRANTY OF
ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT OF THIRD PARTY RIGHTS.
IN NO EVENT SHALL THE COPYRIGHT HOLDER OR HOLDERS INCLUDED IN THIS
NOTICE BE LIABLE FOR ANY CLAIM, OR ANY SPECIAL INDIRECT OR CONSEQUENTIAL
DAMAGES, OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE,
DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER
TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
PERFORMANCE OF THE DATA FILES OR SOFTWARE.
Except as contained in this notice, the name of a copyright holder
shall not be used in advertising or otherwise to promote the sale,
use or other dealings in these Data Files or Software without prior
written authorization of the copyright holder.
```
