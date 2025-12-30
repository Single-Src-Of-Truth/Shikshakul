package com.modulixinfotech.iams.shikshakul.service;

import com.modulixinfotech.iams.shikshakul.domain.model.Role;
import com.modulixinfotech.iams.shikshakul.domain.model.Tenant;
import com.modulixinfotech.iams.shikshakul.domain.model.User;
import com.modulixinfotech.iams.shikshakul.exception.generic.ResourceNotFoundException;
import com.modulixinfotech.iams.shikshakul.repository.UserRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.List;
import java.util.UUID;

import static com.modulixinfotech.iams.shikshakul.domain.model.User.UserStatus.*;

@Service
@RequiredArgsConstructor
@Slf4j
public class UserService {

    private final UserRepository userRepository;
    private final PasswordEncoder passwordEncoder;
    private final AuditService auditService;

    @Transactional
    public User createUser(User user, String rawPassword) {
        if (userRepository.existsByEmailAndTenant(user.getEmail(), user.getTenant())) {
            throw new IllegalArgumentException("User with email already exists in this tenant");
        }

        user.setPasswordHash(passwordEncoder.encode(rawPassword));
        user.setStatus(PENDING);

        User savedUser = userRepository.save(user);
        auditService.logAction(null, user.getTenant(), "USER_CREATED",
                "User", savedUser.getId(), "User created: " + savedUser.getEmail(), true, null);

        log.info("Created user: {} for tenant: {}", savedUser.getEmail(), savedUser.getTenant().getName());
        return savedUser;
    }

    @Transactional(readOnly = true)
    public User getUserById(UUID id) {
        return userRepository.findByIdAndDeletedFalse(id)
                .orElseThrow(() -> new ResourceNotFoundException("User not found with id: " + id));
    }

    @Transactional(readOnly = true)
    public User getUserByEmail(String email, Tenant tenant) {
        return userRepository.findByEmailAndTenantAndDeletedFalse(email, tenant)
                .orElseThrow(() -> new ResourceNotFoundException("User not found with email: " + email));
    }

    @Transactional(readOnly = true)
    public List<User> getUsersByTenant(Tenant tenant) {
        return userRepository.findByTenantAndDeletedFalse(tenant);
    }

    @Transactional(readOnly = true)
    public List<User> getUsersByRole(Role.RoleName roleName, Tenant tenant) {
        return userRepository.findByRoleAndTenant(roleName, tenant);
    }

    @Transactional
    public User updateUser(UUID id, User updatedUser) {
        User user = getUserById(id);

        user.setFirstName(updatedUser.getFirstName());
        user.setLastName(updatedUser.getLastName());
        user.setMobile(updatedUser.getMobile());

        User saved = userRepository.save(user);
        auditService.logAction(user, user.getTenant(), "USER_UPDATED",
                "User", saved.getId(), "User updated", true, null);

        log.info("Updated user: {}", user.getId());
        return saved;
    }

    @Transactional
    public User updateUserStatus(UUID id, User.UserStatus status) {
        User user = getUserById(id);
        user.setStatus(status);

        User saved = userRepository.save(user);
        auditService.logAction(user, user.getTenant(), "USER_STATUS_CHANGED",
                "User", saved.getId(), "Status changed to: " + status, true, null);

        log.info("Updated user status: {} to {}", user.getId(), status);
        return saved;
    }

    @Transactional
    public User addRoleToUser(UUID userId, Role role) {
        User user = getUserById(userId);
        user.getRoles().add(role);

        User saved = userRepository.save(user);
        auditService.logAction(user, user.getTenant(), "ROLE_ADDED",
                "User", saved.getId(), "Role added: " + role.getName(), true, null);

        log.info("Added role {} to user: {}", role.getName(), user.getId());
        return saved;
    }

    @Transactional
    public User removeRoleFromUser(UUID userId, Role role) {
        User user = getUserById(userId);
        user.getRoles().remove(role);

        User saved = userRepository.save(user);
        auditService.logAction(user, user.getTenant(), "ROLE_REMOVED",
                "User", saved.getId(), "Role removed: " + role.getName(), true, null);

        log.info("Removed role {} from user: {}", role.getName(), user.getId());
        return saved;
    }

    @Transactional
    public void changePassword(UUID userId, String oldPassword, String newPassword) {
        User user = getUserById(userId);

        if (!passwordEncoder.matches(oldPassword, user.getPasswordHash())) {
            throw new IllegalArgumentException("Invalid old password");
        }

        user.setPasswordHash(passwordEncoder.encode(newPassword));
        userRepository.save(user);

        auditService.logAction(user, user.getTenant(), "PASSWORD_CHANGED",
                "User", user.getId(), "Password changed", true, null);

        log.info("Password changed for user: {}", user.getId());
    }

    @Transactional
    public void recordLoginAttempt(String email, Tenant tenant, boolean success) {
        userRepository.findByEmailAndTenantAndDeletedFalse(email, tenant)
                .ifPresent(user -> {
                    if (success) {
                        user.setFailedLoginAttempts(0);
                        user.setLockedUntil(null);
                        user.setLastLoginAt(LocalDateTime.now());
                    } else {
                        user.setFailedLoginAttempts(user.getFailedLoginAttempts() + 1);

                        if (user.getFailedLoginAttempts() >= 5) {
                            user.setStatus(LOCKED);
                            user.setLockedUntil(LocalDateTime.now().plusHours(1));
                            log.warn("User account locked due to failed login attempts: {}", email);
                        }
                    }
                    userRepository.save(user);
                });
    }

    @Transactional
    public void unlockUser(UUID userId) {
        User user = getUserById(userId);
        user.setStatus(ACTIVE);
        user.setFailedLoginAttempts(0);
        user.setLockedUntil(null);

        userRepository.save(user);
        auditService.logAction(user, user.getTenant(), "USER_UNLOCKED",
                "User", user.getId(), "User unlocked", true, null);

        log.info("Unlocked user: {}", user.getId());
    }

    @Transactional
    public void deleteUser(UUID id) {
        User user = getUserById(id);
        user.setDeleted(true);
        userRepository.save(user);

        auditService.logAction(user, user.getTenant(), "USER_DELETED",
                "User", user.getId(), "User soft deleted", true, null);

        log.info("Soft deleted user: {}", id);
    }
}