package com.modulixinfotech.iams.shikshakul.config;

import com.modulixinfotech.iams.shikshakul.service.RoleService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.boot.CommandLineRunner;
import org.springframework.stereotype.Component;

@Component
@RequiredArgsConstructor
@Slf4j
public class DataInitializer implements CommandLineRunner {

    private final RoleService roleService;

    @Override
    public void run(String... args) {
        log.info("Initializing application data...");
        roleService.initializeRoles();
        log.info("Application data initialization complete");
    }
}