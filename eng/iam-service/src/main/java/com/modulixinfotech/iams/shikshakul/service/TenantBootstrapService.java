package com.modulixinfotech.iams.shikshakul.service;

import com.modulixinfotech.iams.shikshakul.domain.model.Role;
import com.modulixinfotech.iams.shikshakul.domain.model.Tenant;
import com.modulixinfotech.iams.shikshakul.domain.model.User;
import com.modulixinfotech.iams.shikshakul.repository.UserRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.HashSet;
import java.util.Set;

import static com.modulixinfotech.iams.shikshakul.mapper.request.AuthRequest.CreateTenantRequest;

@Service
@RequiredArgsConstructor
@Slf4j
public class TenantBootstrapService {

    private final TenantService tenantService;
    private final RoleService roleService;
    private final UserRepository userRepository;
    private final PasswordEncoder passwordEncoder;
    private final AuditService auditService;

    @Transactional
    public Tenant createTenantWithAdmin(CreateTenantRequest request) {
        Tenant tenant = Tenant.builder()
                .name(request.getName())
                .domain(request.getDomain())
                .description(request.getDescription())
                .status(Tenant.TenantStatus.ACTIVE)
                .contactEmail(request.getContactEmail())
                .contactPhone(request.getContactPhone())
                .build();

        Tenant savedTenant = tenantService.createTenant(tenant);

        Role adminRole = roleService.getRoleByName(Role.RoleName.ADMIN);
        Set<Role> roles = new HashSet<>();
        roles.add(adminRole);

        User admin = User.builder()
                .tenant(savedTenant)
                .firstName(request.getAdminFirstName())
                .lastName(request.getAdminLastName())
                .email(request.getAdminEmail())
                .passwordHash(passwordEncoder.encode(request.getAdminPassword()))
                .status(User.UserStatus.ACTIVE)
                .emailVerified(true)
                .mobileVerified(false)
                .failedLoginAttempts(0)
                .roles(roles)
                .build();

        User savedAdmin = userRepository.save(admin);

        auditService.logAction(savedAdmin, savedTenant, "TENANT_CREATED",
                "Tenant", savedTenant.getId(),
                "Tenant created with admin user", true, null);

        log.info("Created tenant: {} with admin: {}", savedTenant.getDomain(), savedAdmin.getEmail());

        return savedTenant;
    }
}