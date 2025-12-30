package com.modulixinfotech.iams.shikshakul.domain.model;

import jakarta.persistence.*;
import lombok.*;

@Entity
@Table(name = "roles", indexes = {
        @Index(name = "idx_role_name", columnList = "name")
})
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class Role extends BaseEntity {

    @Enumerated(EnumType.STRING)
    @Column(nullable = false, unique = true, length = 20)
    private RoleName name;

    @Column(length = 255)
    private String description;

    public enum RoleName {
        DEVELOPER,
        ADMIN,
        MANAGEMENT,
        TEACHER,
        STUDENT,
        GUARDIAN
    }
}