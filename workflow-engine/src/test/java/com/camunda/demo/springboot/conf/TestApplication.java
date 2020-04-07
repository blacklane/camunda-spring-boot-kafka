package com.camunda.demo.springboot.conf;

import org.springframework.boot.autoconfigure.EnableAutoConfiguration;
import org.springframework.boot.test.context.TestConfiguration;
import org.springframework.context.annotation.ComponentScan;

import com.camunda.demo.springboot.Application;

@ComponentScan(basePackageClasses={Application.class}) 
@EnableAutoConfiguration()
@TestConfiguration
public class TestApplication {

  // ToDo: mock kafka
}
