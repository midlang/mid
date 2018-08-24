#!/bin/bash

function install_mid() {
	local _prefix=$1
	if [[ -z "$_prefix" ]]; then
		$_prefix=/usr/local
	fi
	cp bin/* $_prefix/bin/
	cp midconfig $_prefix/etc/
	mkdir -p $HOME/.mid
	cp -r templates $HOME/.mid/
	cp -r extentions $HOME/.mid/
}

install_mid $PREFIX
