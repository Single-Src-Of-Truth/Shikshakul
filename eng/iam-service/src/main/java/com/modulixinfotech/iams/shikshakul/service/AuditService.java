package com.modulixinfotech.iams.shikshakul.service;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.modulixinfotech.iams.shikshakul.domain.model.AuditLog;
import com.modulixinfotech.iams.shikshakul.domain.model.Tenant;
import com.modulixinfotech.iams.shikshakul.domain.model.User;
import com.modulixinfotech.iams.shikshakul.repository.AuditLogRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.scheduling.annotation.Async;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Map;
import java.util.UUID;

@Service
@RequiredArgsConstructor
@Slf4j
public class AuditService {

    private final AuditLogRepository auditLogRepository;
    private final ObjectMapper objectMapper;

    @Async
    @Transactional
    public void logAction(User user, Tenant tenant, String action, String entityType,
                          UUID entityId, String description, boolean success, String errorMessage) {
        try {
            AuditLog auditLog = AuditLog.builder()
                    .user(user)
                    .tenant(tenant)
                    .action(action)
                    .entityType(entityType)
                    .entityId(entityId)
                    .metadata(description)
                    .success(success)
                    .errorMessage(errorMessage)
                    .createdAt(LocalDateTime.now())
                    .build();

            auditLogRepository.save(auditLog);
        } catch (Exception e) {
            log.error("Failed to create audit log", e);
        }
    }

    @Async
    @Transactional
    public void logActionWithMetadata(User user, Tenant tenant, String action, String entityType,
                                      UUID entityId, Map<String, Object> metadata, boolean success, String errorMessage) {
        try {
            String metadataJson = objectMapper.writeValueAsString(metadata);

            AuditLog auditLog = AuditLog.builder()
                    .user(user)
                    .tenant(tenant)
                    .action(action)
                    .entityType(entityType)
                    .entityId(entityId)
                    .metadata(metadataJson)
                    .success(success)
                    .errorMessage(errorMessage)
                    .createdAt(LocalDateTime.now())
                    .build();

            auditLogRepository.save(auditLog);
        } catch (JsonProcessingException e) {
            log.error("Failed to serialize audit metadata", e);
        } catch (Exception e) {
            log.error("Failed to create audit log", e);
        }
    }

    @Transactional(readOnly = true)
    public Page<AuditLog> getAuditLogsByUser(User user, Pageable pageable) {
        return auditLogRepository.findByUser(user, pageable);
    }

    @Transactional(readOnly = true)
    public Page<AuditLog> getAuditLogsByTenant(Tenant tenant, Pageable pageable) {
        return auditLogRepository.findByTenant(tenant, pageable);
    }

    @Transactional(readOnly = true)
    public List<AuditLog> getAuditLogsByUserAndDateRange(User user, LocalDateTime start, LocalDateTime end) {
        return auditLogRepository.findByUserAndCreatedAtBetween(user, start, end);
    }

    @Transactional(readOnly = true)
    public List<AuditLog> getAuditLogsByTenantAndDateRange(Tenant tenant, LocalDateTime start, LocalDateTime end) {
        return auditLogRepository.findByTenantAndCreatedAtBetween(tenant, start, end);
    }
}