package com.modulixinfotech.iams.shikshakul.controller;

import com.modulixinfotech.iams.shikshakul.config.AppProperties;
import com.modulixinfotech.iams.shikshakul.domain.model.Tenant;
import com.modulixinfotech.iams.shikshakul.mapper.TenantDTO;
import com.modulixinfotech.iams.shikshakul.mapper.common.ApiResponse;
import com.modulixinfotech.iams.shikshakul.service.TenantBootstrapService;
import com.modulixinfotech.iams.shikshakul.service.TenantService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.UUID;
import java.util.stream.Collectors;

import static com.modulixinfotech.iams.shikshakul.mapper.request.AuthRequest.CreateTenantRequest;

@RestController
@RequestMapping("/api/tenants")
@RequiredArgsConstructor
@Slf4j
public class TenantController {

    private final TenantService tenantService;
    private final TenantBootstrapService tenantBootstrapService;
    private final AppProperties appProperties;

    @PostMapping
    public ResponseEntity<ApiResponse<TenantDTO>> createTenant(
            @Valid @RequestBody CreateTenantRequest request,
            @RequestHeader(value = "MDX-Secret", required = false) String developerSecret) {

        if (developerSecret == null ||
            !developerSecret.equals(appProperties.getSecurity().getDeveloperSecret())) {
            return ResponseEntity
                    .status(HttpStatus.FORBIDDEN)
                    .body(ApiResponse.error("Unauthorized", "Invalid mdx secret"));
        }

        log.info("Creating tenant: {}", request.getDomain());

        Tenant tenant = tenantBootstrapService.createTenantWithAdmin(request);
        TenantDTO tenantDTO = TenantDTO.fromEntity(tenant);

        return ResponseEntity
                .status(HttpStatus.CREATED)
                .body(ApiResponse.success("Tenant created successfully", tenantDTO));
    }

    @GetMapping("/{id}")
    public ResponseEntity<ApiResponse<TenantDTO>> getTenant(@PathVariable UUID id) {
        log.info("Fetching tenant: {}", id);

        Tenant tenant = tenantService.getTenantById(id);
        TenantDTO tenantDTO = TenantDTO.fromEntity(tenant);

        return ResponseEntity.ok(ApiResponse.success(tenantDTO));
    }

    @GetMapping("/domain/{domain}")
    public ResponseEntity<ApiResponse<TenantDTO>> getTenantByDomain(@PathVariable String domain) {
        log.info("Fetching tenant by domain: {}", domain);

        Tenant tenant = tenantService.getTenantByDomain(domain);
        TenantDTO tenantDTO = TenantDTO.fromEntity(tenant);

        return ResponseEntity.ok(ApiResponse.success(tenantDTO));
    }

    @GetMapping
    public ResponseEntity<ApiResponse<List<TenantDTO>>> getAllTenants(
            @RequestHeader(value = "MDX-Secret", required = false) String developerSecret) {

        if (developerSecret == null ||
            !developerSecret.equals(appProperties.getSecurity().getDeveloperSecret())) {
            return ResponseEntity
                    .status(HttpStatus.FORBIDDEN)
                    .body(ApiResponse.error("Unauthorized", "Invalid developer secret"));
        }

        log.info("Fetching all tenants");

        List<TenantDTO> tenants = tenantService.getAllTenants()
                .stream()
                .map(TenantDTO::fromEntity)
                .collect(Collectors.toList());

        return ResponseEntity.ok(ApiResponse.success(tenants));
    }
}