package com.modulixinfotech.iams.shikshakul.repository;

import com.modulixinfotech.iams.shikshakul.domain.model.OAuth2Client;
import com.modulixinfotech.iams.shikshakul.domain.model.Tenant;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Repository
public interface OAuth2ClientRepository extends JpaRepository<OAuth2Client, UUID> {
    Optional<OAuth2Client> findByClientId(String clientId);

    Optional<OAuth2Client> findByClientIdAndDeletedFalse(String clientId);

    List<OAuth2Client> findByTenantAndDeletedFalse(Tenant tenant);

    boolean existsByClientId(String clientId);
}