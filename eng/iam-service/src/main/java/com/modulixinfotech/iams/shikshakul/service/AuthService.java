package com.modulixinfotech.iams.shikshakul.service;

import com.modulixinfotech.iams.shikshakul.domain.model.Invite;
import com.modulixinfotech.iams.shikshakul.domain.model.Role;
import com.modulixinfotech.iams.shikshakul.domain.model.User;
import com.modulixinfotech.iams.shikshakul.repository.UserRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.HashSet;
import java.util.Set;

import static com.modulixinfotech.iams.shikshakul.mapper.request.AuthRequest.RegisterRequest;
import static com.modulixinfotech.iams.shikshakul.mapper.response.AuthResponse.RegisterResponse;

@Service
@RequiredArgsConstructor
@Slf4j
public class AuthService {

    private final UserRepository userRepository;
    private final InviteService inviteService;
    private final PasswordEncoder passwordEncoder;
    private final AuditService auditService;

    @Transactional
    public RegisterResponse registerUser(RegisterRequest request) {
        if (!request.getPassword().equals(request.getConfirmPassword())) {
            throw new IllegalArgumentException("Passwords do not match");
        }

        Invite invite = inviteService.validateInviteToken(request.getInviteToken());

        if (!invite.getEmail().equalsIgnoreCase(request.getEmail())) {
            throw new IllegalArgumentException("Email does not match invite");
        }

        if (userRepository.existsByEmailAndTenant(request.getEmail(), invite.getTenant())) {
            throw new IllegalArgumentException("User already exists with this email in the tenant");
        }

        Set<Role> roles = new HashSet<>();
        roles.add(invite.getRole());

        User user = User.builder()
                .tenant(invite.getTenant())
                .firstName(request.getFirstName())
                .lastName(request.getLastName())
                .email(request.getEmail())
                .mobile(request.getMobile())
                .passwordHash(passwordEncoder.encode(request.getPassword()))
                .status(User.UserStatus.PENDING)
                .emailVerified(false)
                .mobileVerified(false)
                .failedLoginAttempts(0)
                .roles(roles)
                .build();

        User savedUser = userRepository.save(user);

        inviteService.acceptInvite(request.getInviteToken(), savedUser);

        auditService.logAction(savedUser, savedUser.getTenant(), "USER_REGISTERED",
                "User", savedUser.getId(), "User registered via invite", true, null);

        log.info("User registered successfully: {} for tenant: {}",
                savedUser.getEmail(), savedUser.getTenant().getName());

        return RegisterResponse.builder()
                .userId(savedUser.getId())
                .email(savedUser.getEmail())
                .firstName(savedUser.getFirstName())
                .lastName(savedUser.getLastName())
                .status(savedUser.getStatus().name())
                .message("Registration successful. Your account is pending admin approval.")
                .build();
    }

    @Transactional
    public void activateUser(User user) {
        user.setStatus(User.UserStatus.ACTIVE);
        userRepository.save(user);

        auditService.logAction(user, user.getTenant(), "USER_ACTIVATED",
                "User", user.getId(), "User account activated", true, null);

        log.info("User activated: {}", user.getEmail());
    }

    @Transactional
    public void verifyEmail(User user) {
        user.setEmailVerified(true);
        userRepository.save(user);

        auditService.logAction(user, user.getTenant(), "EMAIL_VERIFIED",
                "User", user.getId(), "Email verified", true, null);

        log.info("Email verified for user: {}", user.getEmail());
    }

    @Transactional
    public void verifyMobile(User user) {
        user.setMobileVerified(true);
        userRepository.save(user);

        auditService.logAction(user, user.getTenant(), "MOBILE_VERIFIED",
                "User", user.getId(), "Mobile verified", true, null);

        log.info("Mobile verified for user: {}", user.getEmail());
    }
}