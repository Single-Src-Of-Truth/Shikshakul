package com.modulixinfotech.iams.shikshakul.domain.model;

import jakarta.persistence.*;
import lombok.*;

import java.util.Set;

@Entity
@Table(name = "oauth2_clients", indexes = {
        @Index(name = "idx_oauth_client_id", columnList = "clientId"),
        @Index(name = "idx_oauth_tenant", columnList = "tenant_id")
})
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class OAuth2Client extends BaseEntity {

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "tenant_id", nullable = false)
    private Tenant tenant;

    @Column(nullable = false, unique = true, length = 100)
    private String clientId;

    @Column(length = 255)
    private String clientSecret;

    @Column(nullable = false, length = 200)
    private String clientName;

    @ElementCollection(fetch = FetchType.EAGER)
    @CollectionTable(name = "oauth2_client_redirect_uris",
            joinColumns = @JoinColumn(name = "client_id"))
    @Column(name = "redirect_uri", length = 500)
    private Set<String> redirectUris;

    @ElementCollection(fetch = FetchType.EAGER)
    @CollectionTable(name = "oauth2_client_grant_types",
            joinColumns = @JoinColumn(name = "client_id"))
    @Column(name = "grant_type", length = 50)
    private Set<String> grantTypes;

    @ElementCollection(fetch = FetchType.EAGER)
    @CollectionTable(name = "oauth2_client_scopes",
            joinColumns = @JoinColumn(name = "client_id"))
    @Column(name = "scope", length = 100)
    private Set<String> scopes;

    @Column(nullable = false)
    private Boolean requirePkce = true;

    @Column(nullable = false)
    private Boolean requireAuthorizationConsent = false;

    @Column(nullable = false)
    private Integer accessTokenTtl = 900;

    @Column(nullable = false)
    private Integer refreshTokenTtl = 2592000;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false, length = 20)
    private ClientStatus status;

    public enum ClientStatus {
        ACTIVE,
        INACTIVE,
        SUSPENDED
    }
}