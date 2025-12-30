package com.modulixinfotech.iams.shikshakul.mapper;

import com.modulixinfotech.iams.shikshakul.mapper.common.BaseDTO;
import com.modulixinfotech.iams.shikshakul.domain.model.User;
import com.modulixinfotech.iams.shikshakul.domain.model.User.UserStatus;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.NoArgsConstructor;
import lombok.experimental.SuperBuilder;

import java.time.LocalDateTime;
import java.util.Set;
import java.util.UUID;
import java.util.stream.Collectors;

@Data
@SuperBuilder
@NoArgsConstructor
@AllArgsConstructor
@EqualsAndHashCode(callSuper = true)
public class UserDTO extends BaseDTO {
    private UUID tenantId;
    private String tenantName;
    private String firstName;
    private String lastName;
    private String email;
    private String mobile;
    private UserStatus status;
    private Boolean emailVerified;
    private Boolean mobileVerified;
    private LocalDateTime lastLoginAt;
    private Set<String> roles;

    public static UserDTO fromEntity(User user) {
        if (user == null) return null;

        return UserDTO.builder()
                .id(user.getId())
                .tenantId(user.getTenant().getId())
                .tenantName(user.getTenant().getName())
                .firstName(user.getFirstName())
                .lastName(user.getLastName())
                .email(user.getEmail())
                .mobile(user.getMobile())
                .status(user.getStatus())
                .emailVerified(user.getEmailVerified())
                .mobileVerified(user.getMobileVerified())
                .lastLoginAt(user.getLastLoginAt())
                .roles(user.getRoles().stream()
                        .map(role -> role.getName().name())
                        .collect(Collectors.toSet()))
                .createdAt(user.getCreatedAt())
                .updatedAt(user.getUpdatedAt())
                .build();
    }
}
