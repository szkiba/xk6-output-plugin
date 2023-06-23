#!/usr/bin/env node

import { serve } from 'xk6-output-plugin-js'

class Example {
  init (params) {
    console.info('init')
    return { description: 'example-js plugin' }
  }

  start () {
    console.info('start')
  }

  stop () {
    console.info('stop')
  }

  addMetrics (metrics) {
    console.info('metrics')
    metrics.forEach(metric => {
      console.info(
        { 'metric.type': metric.type, 'metric.contains': metric.contains },
        metric.name
      )
    })
  }

  addSamples (samples) {
    console.info('samples')
    samples.forEach(sample => {
      console.info(
        { 'sample.time': new Date(sample.time).toISOString(), 'sample.value': sample.value },
        sample.metric
      )
    })
  }
}

serve(new Example())
