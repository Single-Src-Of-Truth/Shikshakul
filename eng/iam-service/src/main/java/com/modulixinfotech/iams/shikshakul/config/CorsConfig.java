package com.modulixinfotech.iams.shikshakul.config;

import lombok.RequiredArgsConstructor;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.cors.CorsConfiguration;
import org.springframework.web.cors.CorsConfigurationSource;
import org.springframework.web.cors.UrlBasedCorsConfigurationSource;

import java.util.Arrays;
import java.util.List;

@Configuration
@RequiredArgsConstructor
public class CorsConfig {

    private final AppProperties appProperties;

    @Bean
    public CorsConfigurationSource corsConfigurationSource() {
        CorsConfiguration configuration = new CorsConfiguration();

        String[] origins = appProperties.getCors().getAllowedOrigins().split(",");
        configuration.setAllowedOrigins(Arrays.asList(origins));

        String[] methods = appProperties.getCors().getAllowedMethods().split(",");
        configuration.setAllowedMethods(Arrays.asList(methods));

        if ("*".equals(appProperties.getCors().getAllowedHeaders())) {
            configuration.addAllowedHeader("*");
        } else {
            String[] headers = appProperties.getCors().getAllowedHeaders().split(",");
            configuration.setAllowedHeaders(Arrays.asList(headers));
        }

        configuration.setAllowCredentials(appProperties.getCors().getAllowCredentials());
        configuration.setMaxAge(appProperties.getCors().getMaxAge());

        configuration.setExposedHeaders(List.of("Authorization", "Set-Cookie"));

        UrlBasedCorsConfigurationSource source = new UrlBasedCorsConfigurationSource();
        source.registerCorsConfiguration("/**", configuration);
        return source;
    }
}