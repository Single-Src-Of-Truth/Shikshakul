package com.modulixinfotech.iams.shikshakul.config;

import lombok.Data;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.stereotype.Component;

@Data
@Component
@ConfigurationProperties(prefix = "app")
public class AppProperties {

    private Security security = new Security();
    private Oauth2 oauth2 = new Oauth2();
    private Invite invite = new Invite();
    private Cookie cookie = new Cookie();
    private Cors cors = new Cors();

    @Data
    public static class Security {
        private String developerSecret;
    }

    @Data
    public static class Oauth2 {
        private String issuerUrl;
        private Integer authorizationCodeTtl = 300;
        private Integer accessTokenTtl = 900;
        private Integer refreshTokenTtl = 2592000;
    }

    @Data
    public static class Invite {
        private Integer defaultTtlHours = 48;
    }

    @Data
    public static class Cookie {
        private String name = "mdx_refresh";
        private String domain = "localhost";
        private Boolean secure = false;
        private String sameSite = "Lax";
        private Integer maxAge = 2592000;
    }

    @Data
    public static class Cors {
        private String allowedOrigins;
        private String allowedMethods;
        private String allowedHeaders;
        private Boolean allowCredentials;
        private Long maxAge;
    }
}