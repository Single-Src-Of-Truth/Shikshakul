package com.modulixinfotech.iams.shikshakul;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class App {

    public static void main(String[] args) {
        SpringApplication.run(App.class, args);
        System.out.println("\n\n ** Identity and Access Management Service is up and running 👍... ** \n\n");
    }

}
