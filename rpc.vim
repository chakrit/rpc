" RPC minimal syntax file for vim
" Language: RPC
" Maintainer: Chakrit Wichian <service@chakrit.net>

if exists("b:current_syntax")
  finish
endif

" see lexer/keywords.go
syn keyword rpcDataTypes bool data double float include int list long map string time
hi def link rpcDataTypes Statement

syn keyword rpcBlocks root rpc type namespace option
hi def link rpcBlocks Type

syn match rpcComment '//.*$'
hi def link rpcComment Comment

syn match rpcNum '\d\+'
syn match rpcStr '"[^"]*"'
hi def link rpcNum Constant
hi def link rpcStr Constant
