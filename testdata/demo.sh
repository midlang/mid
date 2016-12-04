#!/bin/bash

kinds='
default
beans
'

for kind in $kinds
do
midc \
	-I ./demo.mid \
	-Ogo=generated/go_$kind \
	-Ocpp=generated/cpp_$kind \
	-Eautogen_decl="// NOTE: generated file, DON'T edit!!" \
	-Ecpp:unordered_map \
	-Xproto \
	-K $kind \
	--log=info
done
