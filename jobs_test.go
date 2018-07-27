package rundeck_test

import (
	rundeck "github.com/andrewmeissner/go-rundeck"
)

var jobDef1 = `- defaultTab: output
  description: ''
  executionEnabled: true
  loglevel: INFO
  name: who_and_where
  nodeFilterEditable: false
  schedule:
    dayofmonth:
      day: 1/1
    month: '*'
    time:
      hour: '*'
      minute: 0/1
      seconds: '0'
    year: '*'
  scheduleEnabled: true
  sequence:
    commands:
    - exec: whoami
    - exec: pwd
    keepgoing: false
    strategy: node-first
  timeZone: MST7MDT`

var jobDef2 = `- defaultTab: output
  description: ''
  executionEnabled: true
  loglevel: INFO
  name: hostname
  nodeFilterEditable: false
  schedule:
    dayofmonth:
      day: 1/1
    month: '*'
    time:
      hour: '*'
      minute: 0/1
      seconds: '0'
    year: '*'
  scheduleEnabled: true
  sequence:
    commands:
    - exec: hostname
    keepgoing: false
    strategy: node-first
  timeZone: MST7MDT`

func standup(client *rundeck.Client) {

}
