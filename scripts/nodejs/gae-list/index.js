#! /usr/bin/env node
'use strict';

import Table from "cli-table";
import yargs from "yargs";
import {hideBin} from "yargs/helpers";
import chalk from "chalk";
import shell from "shelljs"
import inquirer from "inquirer";
import acp from 'inquirer-autocomplete-prompt';

import lib from "../common-lib/lib.js";

const AutocompletePrompt = acp


if (!shell.which('gcloud')) {
    console.log(chalk.yellow.bold('Sorry, this script requires gcloud'));
    console.log(chalk.red.bold('Aborting...'));
    shell.exit(1);
  }

const argv = yargs(hideBin(process.argv))
    .version('0.0.1')
    .usage('Usage: $ $0')
    .wrap(yargs.terminalWidth)
    .option('p', {
        alias: 'project',
        description: 'Use different project then current context.',
        default: lib.gcpConfig.core.project
    })
    .option('s', {
        description: 'Optionally specify the service to query.',
        alias: 'service'
    })
    .option('f', {
        description: 'Apply a filter',
        alias: 'filter',
        default: 'version.servingStatus=SERVING'
    })
    .example('$ $0 -s shipping-service\nOR\n$ $0 --project np-sam-project')
    .help('h')
    .alias('h', 'help')
    .argv;

var data = "";
var cmd = `gcloud app versions list --project ${argv.p} --format json`
if(argv.s) {
    cmd = cmd + ` --service ${argv.s}`
}
if(argv.f) {
    cmd = cmd + ` --filter="${argv.f}"`
}
// console.log(cmd)
data = JSON.parse(shell.exec(cmd, {silent:true}).stdout)


var table = new Table({
    head: ['SERVICE', 'VERSION', 'TRAFFIC SPLIT'],
    chars: { 'top': '═' , 'top-mid': '╤' , 'top-left': '╔' , 'top-right': '╗'
         , 'bottom': '═' , 'bottom-mid': '╧' , 'bottom-left': '╚' , 'bottom-right': '╝'
         , 'left': '║' , 'left-mid': '╟' , 'mid': '─' , 'mid-mid': '┼'
         , 'right': '║' , 'right-mid': '╢' , 'middle': '│' },
    style: {
        compact : true, 
        'padding-left' : 1
    }
})

data.forEach(element => {
    var row = [element.service, element.version.id, `${(element.traffic_split*100)}%`]
    table.push(row)
});

console.log(table.toString());
