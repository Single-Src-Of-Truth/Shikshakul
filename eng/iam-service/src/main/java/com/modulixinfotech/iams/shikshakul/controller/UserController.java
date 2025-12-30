package com.modulixinfotech.iams.shikshakul.controller;

import com.modulixinfotech.iams.shikshakul.domain.model.Tenant;
import com.modulixinfotech.iams.shikshakul.domain.model.User;
import com.modulixinfotech.iams.shikshakul.mapper.UserDTO;
import com.modulixinfotech.iams.shikshakul.mapper.common.ApiResponse;
import com.modulixinfotech.iams.shikshakul.security.userdetails.CustomUserDetails;
import com.modulixinfotech.iams.shikshakul.service.TenantService;
import com.modulixinfotech.iams.shikshakul.service.UserService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.UUID;
import java.util.stream.Collectors;

import static com.modulixinfotech.iams.shikshakul.mapper.request.AuthRequest.ChangePasswordRequest;

@RestController
@RequestMapping("/api/users")
@RequiredArgsConstructor
@Slf4j
public class UserController {

    private final UserService userService;
    private final TenantService tenantService;

    @GetMapping("/me")
    public ResponseEntity<ApiResponse<UserDTO>> getCurrentUser(
            @AuthenticationPrincipal CustomUserDetails userDetails) {
        log.info("Fetching current user info");

        User user = userService.getUserById(userDetails.getUser().getId());
        UserDTO userDTO = UserDTO.fromEntity(user);

        return ResponseEntity.ok(ApiResponse.success(userDTO));
    }

    @GetMapping("/{id}")
    public ResponseEntity<ApiResponse<UserDTO>> getUserById(@PathVariable UUID id) {
        log.info("Fetching user: {}", id);

        User user = userService.getUserById(id);
        UserDTO userDTO = UserDTO.fromEntity(user);

        return ResponseEntity.ok(ApiResponse.success(userDTO));
    }

    @GetMapping("/tenant/{tenantId}")
    public ResponseEntity<ApiResponse<List<UserDTO>>> getUsersByTenant(@PathVariable UUID tenantId) {
        log.info("Fetching users for tenant: {}", tenantId);

        Tenant tenant = tenantService.getTenantById(tenantId);
        List<UserDTO> users = userService.getUsersByTenant(tenant)
                .stream()
                .map(UserDTO::fromEntity)
                .collect(Collectors.toList());

        return ResponseEntity.ok(ApiResponse.success(users));
    }

    @PutMapping("/{id}/activate")
    public ResponseEntity<ApiResponse<UserDTO>> activateUser(@PathVariable UUID id) {
        log.info("Activating user: {}", id);

        User user = userService.updateUserStatus(id, User.UserStatus.ACTIVE);
        UserDTO userDTO = UserDTO.fromEntity(user);

        return ResponseEntity.ok(ApiResponse.success("User activated successfully", userDTO));
    }

    @PutMapping("/{id}/suspend")
    public ResponseEntity<ApiResponse<UserDTO>> suspendUser(@PathVariable UUID id) {
        log.info("Suspending user: {}", id);

        User user = userService.updateUserStatus(id, User.UserStatus.SUSPENDED);
        UserDTO userDTO = UserDTO.fromEntity(user);

        return ResponseEntity.ok(ApiResponse.success("User suspended successfully", userDTO));
    }

    @PutMapping("/{id}/unlock")
    public ResponseEntity<ApiResponse<UserDTO>> unlockUser(@PathVariable UUID id) {
        log.info("Unlocking user: {}", id);

        userService.unlockUser(id);
        User user = userService.getUserById(id);
        UserDTO userDTO = UserDTO.fromEntity(user);

        return ResponseEntity.ok(ApiResponse.success("User unlocked successfully", userDTO));
    }

    @PostMapping("/change-password")
    public ResponseEntity<ApiResponse<Void>> changePassword(
            @Valid @RequestBody ChangePasswordRequest request,
            @AuthenticationPrincipal CustomUserDetails userDetails) {

        if (!request.getNewPassword().equals(request.getConfirmPassword())) {
            return ResponseEntity
                    .badRequest()
                    .body(ApiResponse.error("Passwords do not match"));
        }

        log.info("Changing password for user: {}", userDetails.getUsername());

        userService.changePassword(
                userDetails.getUser().getId(),
                request.getCurrentPassword(),
                request.getNewPassword()
        );

        return ResponseEntity.ok(ApiResponse.success("Password changed successfully", null));
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<ApiResponse<Void>> deleteUser(@PathVariable UUID id) {
        log.info("Deleting user: {}", id);

        userService.deleteUser(id);

        return ResponseEntity.ok(ApiResponse.success("User deleted successfully", null));
    }
}