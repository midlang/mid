#!/bin/bash

midc \
	-I ./demo.mid \
	-Ogo=generated/go \
	-Ocpp=generated/cpp \
	-Eautogen_decl="// NOTE: generated file, DON'T edit!!" \
	-Ecpp:unordered_map \
	-Xproto \
	--log=debug
