'use strict';

import ini from "ini";
import shell from "shelljs"
import fuzzy from "fuzzy"
import inquirer from "inquirer";
import _ from "lodash"


var HOME = shell.env.HOME
const raw_gcpConfig = shell.cat(`${HOME}/.config/gcloud/configurations/config_default`)
const gcpConfig = ini.parse(raw_gcpConfig)
const gcpCacheDir = `${HOME}/.config/gcloud/cache/${gcpConfig.core.account}`
const projectListFile = `${gcpCacheDir}/projects_list`

function refreshProjectList() {
    console.log("Refreshing project list cache!")
    shell.mkdir('-p', gcpCacheDir)
    shell.exec(`gcloud projects list | tail -n +2 | awk '{print $1}' | grep -v 'sys-' | sort > ${projectListFile}`)
}

function getProjectList() {
    if (shell.test('-f', projectListFile)) {
        return shell.cat(projectListFile).split('\n')
    } else {
        refreshProjectList()
        return getProjectList()
    }
}

export default {
    gcpConfig: gcpConfig,
    projectListFile: projectListFile,
    getProjectList: getProjectList(),
    searchProjects: function(answers, input) {
        var projectList = getProjectList()
        input = input || '';
        return new Promise(function (resolve) {
            setTimeout(function () {
            var fuzzyResult = fuzzy.filter(input, projectList);
            const results = fuzzyResult.map(function (el) {
                return el.original;
            });

            results.splice(5, 0, new inquirer.Separator());
            results.push(new inquirer.Separator());
            resolve(results);
            }, _.random(30, 500));
        });
    },
    // Returns clusters as a Map of key cluster with value of the zone
    getProjectClusters: function (project) {
        var rawClusters = shell.exec(`gcloud container clusters list --project ${project}`, {silent:true}).stdout.split('\n').slice(1);
        var clusters = new Map();
        rawClusters.forEach(clusterLine => {
            var splitz = clusterLine.split(/[ ]+/);
            clusters.set(splitz[0], splitz[1])
        })
        return clusters
    }
}