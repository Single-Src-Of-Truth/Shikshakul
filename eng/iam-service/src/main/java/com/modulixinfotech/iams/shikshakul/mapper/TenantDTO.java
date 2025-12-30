package com.modulixinfotech.iams.shikshakul.mapper;

import com.modulixinfotech.iams.shikshakul.mapper.common.BaseDTO;
import com.modulixinfotech.iams.shikshakul.domain.model.Tenant;
import com.modulixinfotech.iams.shikshakul.domain.model.Tenant.TenantStatus;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.NoArgsConstructor;
import lombok.experimental.SuperBuilder;

@Data
@SuperBuilder
@NoArgsConstructor
@AllArgsConstructor
@EqualsAndHashCode(callSuper = true)
public class TenantDTO extends BaseDTO {
    private String name;
    private String domain;
    private String description;
    private TenantStatus status;
    private String contactEmail;
    private String contactPhone;

    public static TenantDTO fromEntity(Tenant tenant) {
        if (tenant == null) return null;

        return TenantDTO.builder()
                .id(tenant.getId())
                .name(tenant.getName())
                .domain(tenant.getDomain())
                .description(tenant.getDescription())
                .status(tenant.getStatus())
                .contactEmail(tenant.getContactEmail())
                .contactPhone(tenant.getContactPhone())
                .createdAt(tenant.getCreatedAt())
                .updatedAt(tenant.getUpdatedAt())
                .build();
    }
}