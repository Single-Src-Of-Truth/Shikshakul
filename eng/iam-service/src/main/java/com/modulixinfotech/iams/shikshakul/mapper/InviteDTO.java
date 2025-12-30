package com.modulixinfotech.iams.shikshakul.mapper;

import com.modulixinfotech.iams.shikshakul.mapper.common.BaseDTO;
import com.modulixinfotech.iams.shikshakul.domain.model.Invite;
import com.modulixinfotech.iams.shikshakul.domain.model.Invite.InviteStatus;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.NoArgsConstructor;
import lombok.experimental.SuperBuilder;

import java.time.LocalDateTime;
import java.util.UUID;

@Data
@SuperBuilder
@NoArgsConstructor
@AllArgsConstructor
@EqualsAndHashCode(callSuper = true)
public class InviteDTO extends BaseDTO {
    private UUID token;
    private UUID tenantId;
    private String tenantName;
    private String email;
    private String mobile;
    private String role;
    private UUID invitedById;
    private String invitedByName;
    private LocalDateTime expiresAt;
    private InviteStatus status;
    private LocalDateTime acceptedAt;

    public static InviteDTO fromEntity(Invite invite) {
        if (invite == null) return null;

        return InviteDTO.builder()
                .id(invite.getId())
                .token(invite.getToken())
                .tenantId(invite.getTenant().getId())
                .tenantName(invite.getTenant().getName())
                .email(invite.getEmail())
                .mobile(invite.getMobile())
                .role(invite.getRole().getName().name())
                .invitedById(invite.getInvitedBy() != null ? invite.getInvitedBy().getId() : null)
                .invitedByName(invite.getInvitedBy() != null ?
                        invite.getInvitedBy().getFirstName() + " " + invite.getInvitedBy().getLastName() : null)
                .expiresAt(invite.getExpiresAt())
                .status(invite.getStatus())
                .acceptedAt(invite.getAcceptedAt())
                .createdAt(invite.getCreatedAt())
                .updatedAt(invite.getUpdatedAt())
                .build();
    }
}