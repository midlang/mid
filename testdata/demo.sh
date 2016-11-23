#!/bin/bash

midc -I ./demo.mid -Ogo=generated/go -Ocpp=generated/cpp --log=debug -Eautogen_decl="// NOTE: generated file, DON'T edit!!" -Ecpp:unordered_map
