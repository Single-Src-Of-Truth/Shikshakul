package com.modulixinfotech.iams.shikshakul.service;

import com.modulixinfotech.iams.shikshakul.domain.model.Tenant;
import com.modulixinfotech.iams.shikshakul.exception.generic.ResourceNotFoundException;
import com.modulixinfotech.iams.shikshakul.repository.TenantRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.UUID;

@Service
@RequiredArgsConstructor
@Slf4j
public class TenantService {

    private final TenantRepository tenantRepository;

    @Transactional
    public Tenant createTenant(Tenant tenant) {
        if (tenantRepository.existsByDomain(tenant.getDomain())) {
            throw new IllegalArgumentException("Tenant with domain already exists: " + tenant.getDomain());
        }

        log.info("Creating new tenant: {}", tenant.getName());
        return tenantRepository.save(tenant);
    }

    @Transactional(readOnly = true)
    public Tenant getTenantById(UUID id) {
        return tenantRepository.findByIdAndDeletedFalse(id)
                .orElseThrow(() -> new ResourceNotFoundException("Tenant not found with id: " + id));
    }

    @Transactional(readOnly = true)
    public Tenant getTenantByDomain(String domain) {
        return tenantRepository.findByDomainAndDeletedFalse(domain)
                .orElseThrow(() -> new ResourceNotFoundException("Tenant not found with domain: " + domain));
    }

    @Transactional(readOnly = true)
    public List<Tenant> getAllTenants() {
        return tenantRepository.findAll();
    }

    @Transactional
    public Tenant updateTenant(UUID id, Tenant updatedTenant) {
        Tenant tenant = getTenantById(id);

        tenant.setName(updatedTenant.getName());
        tenant.setDescription(updatedTenant.getDescription());
        tenant.setStatus(updatedTenant.getStatus());
        tenant.setContactEmail(updatedTenant.getContactEmail());
        tenant.setContactPhone(updatedTenant.getContactPhone());

        log.info("Updating tenant: {}", tenant.getId());
        return tenantRepository.save(tenant);
    }

    @Transactional
    public void deleteTenant(UUID id) {
        Tenant tenant = getTenantById(id);
        tenant.setDeleted(true);
        tenantRepository.save(tenant);
        log.info("Soft deleted tenant: {}", id);
    }
}