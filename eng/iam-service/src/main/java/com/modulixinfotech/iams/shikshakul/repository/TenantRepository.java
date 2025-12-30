package com.modulixinfotech.iams.shikshakul.repository;

import com.modulixinfotech.iams.shikshakul.domain.model.Tenant;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;
import java.util.UUID;

@Repository
public interface TenantRepository extends JpaRepository<Tenant, UUID> {
    Optional<Tenant> findByDomain(String domain);

    Optional<Tenant> findByIdAndDeletedFalse(UUID id);

    Optional<Tenant> findByDomainAndDeletedFalse(String domain);

    boolean existsByDomain(String domain);
}