package com.modulixinfotech.iams.shikshakul.domain.model;

import jakarta.persistence.*;
import lombok.*;

import java.time.LocalDateTime;
import java.util.UUID;

@Entity
@Table(name = "invites", indexes = {
        @Index(name = "idx_invite_token", columnList = "token"),
        @Index(name = "idx_invite_tenant", columnList = "tenant_id"),
        @Index(name = "idx_invite_status", columnList = "status"),
        @Index(name = "idx_invite_expires", columnList = "expiresAt")
})
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class Invite extends BaseEntity {

    @Column(nullable = false, unique = true)
    private UUID token;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "tenant_id", nullable = false)
    private Tenant tenant;

    @Column(nullable = false, length = 100)
    private String email;

    @Column(length = 20)
    private String mobile;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "role_id", nullable = false)
    private Role role;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "invited_by_id")
    private User invitedBy;

    @Column(nullable = false)
    private LocalDateTime expiresAt;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false, length = 20)
    private InviteStatus status;

    @Column
    private LocalDateTime acceptedAt;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "accepted_by_id")
    private User acceptedBy;

    public enum InviteStatus {
        PENDING,
        ACCEPTED,
        EXPIRED,
        REVOKED
    }
}