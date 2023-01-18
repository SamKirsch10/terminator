#! /usr/bin/env node
'use strict';


import yargs from "yargs";
import {hideBin} from "yargs/helpers";
import chalk from "chalk";
import boxen from "boxen";
import shell from "shelljs";
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
    .version('0.0.2')
    .usage('Usage: $ $0')
    .wrap(yargs.terminalWidth)
    .command('[cluster]', 'Optionally specify the cluster to set kubectl context without going thru menus.')
    .example('$ $0 sam-k8s-cluster-us-east-1b\nOR\n$ $0')
    .option('refresh', {
            alias: 'r',
            default: false,
            describe: 'Force refresh GCP project listing',
            type: 'boolean'
        }
    )
    .help('h')
    .alias('h', 'help')
    .parse()
    ;


if ( argv.refresh === true ) {
    shell.rm(lib.projectListFile)
}


const boxenOptions = {
 padding: 1,
 margin: 1,
 borderStyle: "double",
 borderColor: "blue",
};
const intro = chalk.white.bold("K8s Context Switcher")
const msgBox = boxen( intro, boxenOptions );
console.log(msgBox);

if (argv._.length > 0) {
    var cluster_arg = argv._[0]
    var clusters = lib.getProjectClusters(lib.gcpConfig.core.project)
    if (!clusters.has(cluster_arg)) {
        console.log(chalk.red.bold(`The current project [${lib.gcpConfig.core.project}] does not have a cluster named '${cluster_arg}'!`))
        shell.exit(1)
    }
    shell.exec(`gcloud container clusters get-credentials ${cluster_arg} --zone ${clusters.get(cluster_arg)}`)
} else {
    inquirer.registerPrompt('autocomplete', AutocompletePrompt);
    var clusters = new Map();
    var promptFields = [
        {
            type: 'autocomplete',
            name: 'gcpProject',
            default: lib.gcpConfig.core.project,
            message: 'GCP Project',
            emptyText: 'project not found!',
            source: lib.searchProjects,
            suggestOnly: true,
            validate: function (val) {
                return val ? true : 'Type something!';
              },
        },
        {
            type: 'list',
            name: 'cluster',
            choices: function(previousAnswers, input) {
                clusters = lib.getProjectClusters(previousAnswers['gcpProject'])
                return Array.from( clusters.keys() )
            }
        }
    ]
    inquirer
    .prompt(promptFields)
    .then((answers) => {
        shell.exec(`gcloud config set project ${answers['gcpProject']}`)
        shell.exec(`gcloud container clusters get-credentials ${answers['cluster']} --zone ${clusters.get(answers['cluster'])}`)
    });

}