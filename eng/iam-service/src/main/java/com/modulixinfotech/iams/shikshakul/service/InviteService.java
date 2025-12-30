package com.modulixinfotech.iams.shikshakul.service;

import com.modulixinfotech.iams.shikshakul.config.AppProperties;
import com.modulixinfotech.iams.shikshakul.domain.model.Invite;
import com.modulixinfotech.iams.shikshakul.domain.model.Role;
import com.modulixinfotech.iams.shikshakul.domain.model.Tenant;
import com.modulixinfotech.iams.shikshakul.domain.model.User;
import com.modulixinfotech.iams.shikshakul.exception.generic.ResourceNotFoundException;
import com.modulixinfotech.iams.shikshakul.repository.InviteRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.List;
import java.util.UUID;

import static com.modulixinfotech.iams.shikshakul.domain.model.Invite.InviteStatus.PENDING;

@Service
@RequiredArgsConstructor
@Slf4j
public class InviteService {

    private final InviteRepository inviteRepository;
    private final AppProperties appProperties;
    private final AuditService auditService;

    @Transactional
    public Invite createInvite(Tenant tenant, String email, String mobile, Role role, User invitedBy) {
        UUID token = UUID.randomUUID();
        LocalDateTime expiresAt = LocalDateTime.now()
                .plusHours(appProperties.getInvite().getDefaultTtlHours());

        Invite invite = Invite.builder()
                .token(token)
                .tenant(tenant)
                .email(email)
                .mobile(mobile)
                .role(role)
                .invitedBy(invitedBy)
                .expiresAt(expiresAt)
                .status(PENDING)
                .build();

        Invite saved = inviteRepository.save(invite);
        auditService.logAction(invitedBy, tenant, "INVITE_CREATED",
                "Invite", saved.getId(), "Invite sent to: " + email, true, null);

        log.info("Created invite for {} in tenant: {}", email, tenant.getName());
        return saved;
    }

    @Transactional(readOnly = true)
    public Invite getInviteByToken(UUID token) {
        return inviteRepository.findByToken(token)
                .orElseThrow(() -> new ResourceNotFoundException("Invite not found"));
    }

    @Transactional(readOnly = true)
    public Invite validateInviteToken(UUID token) {
        Invite invite = inviteRepository.findByTokenAndStatus(token, PENDING)
                .orElseThrow(() -> new IllegalArgumentException("Invalid or already used invite"));

        if (invite.getExpiresAt().isBefore(LocalDateTime.now())) {
            throw new IllegalArgumentException("Invite has expired");
        }

        return invite;
    }

    @Transactional
    public Invite acceptInvite(UUID token, User acceptedBy) {
        Invite invite = validateInviteToken(token);

        invite.setStatus(Invite.InviteStatus.ACCEPTED);
        invite.setAcceptedAt(LocalDateTime.now());
        invite.setAcceptedBy(acceptedBy);

        Invite saved = inviteRepository.save(invite);
        auditService.logAction(acceptedBy, invite.getTenant(), "INVITE_ACCEPTED",
                "Invite", saved.getId(), "Invite accepted by: " + acceptedBy.getEmail(), true, null);

        log.info("Invite accepted: {}", token);
        return saved;
    }

    @Transactional
    public void revokeInvite(UUID inviteId, User revokedBy) {
        Invite invite = inviteRepository.findById(inviteId)
                .orElseThrow(() -> new ResourceNotFoundException("Invite not found"));

        invite.setStatus(Invite.InviteStatus.REVOKED);
        inviteRepository.save(invite);

        auditService.logAction(revokedBy, invite.getTenant(), "INVITE_REVOKED",
                "Invite", invite.getId(), "Invite revoked", true, null);

        log.info("Invite revoked: {}", inviteId);
    }

    @Transactional(readOnly = true)
    public List<Invite> getInvitesByTenant(Tenant tenant) {
        return inviteRepository.findByTenantAndDeletedFalse(tenant);
    }

    @Transactional(readOnly = true)
    public List<Invite> getPendingInvitesByTenant(Tenant tenant) {
        return inviteRepository.findByTenantAndStatusAndDeletedFalse(tenant, PENDING);
    }

    @Transactional
    public void expireOldInvites() {
        List<Invite> expiredInvites = inviteRepository.findByStatusAndExpiresAtBefore(
                PENDING, LocalDateTime.now());

        for (Invite invite : expiredInvites) {
            invite.setStatus(Invite.InviteStatus.EXPIRED);
            inviteRepository.save(invite);
        }

        if (!expiredInvites.isEmpty()) {
            log.info("Expired {} old invites", expiredInvites.size());
        }
    }
}