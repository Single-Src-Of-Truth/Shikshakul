package com.modulixinfotech.iams.shikshakul.repository;

import com.modulixinfotech.iams.shikshakul.domain.model.Invite;
import com.modulixinfotech.iams.shikshakul.domain.model.Invite.InviteStatus;
import com.modulixinfotech.iams.shikshakul.domain.model.Tenant;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Repository
public interface InviteRepository extends JpaRepository<Invite, UUID> {
    Optional<Invite> findByToken(UUID token);

    Optional<Invite> findByTokenAndStatus(UUID token, InviteStatus status);

    List<Invite> findByTenantAndDeletedFalse(Tenant tenant);

    List<Invite> findByTenantAndStatusAndDeletedFalse(Tenant tenant, InviteStatus status);

    List<Invite> findByStatusAndExpiresAtBefore(InviteStatus status, LocalDateTime dateTime);
}