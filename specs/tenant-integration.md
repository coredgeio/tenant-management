# 1. Introduction
As our platform continues to integrate with multiple third-party services, managing tenant-specific information such as KYC, KYB, payment method configuration, and tenant type has become increasingly complex. This complexity is compounded by varying client requirements, update patterns, and data access mechanisms. 

To address this challenge, we propose a centralized solution for managing tenant-related data that ensures consistency, scalability, and flexibility across clients. This solution introduces a modular architecture with dedicated managers and a shared plugin-based system, allowing seamless support for existing and future clients. 

This document outlines the proposed architecture, configuration strategy, execution flow, and implementation plan to realize this centralized tenant information management system. 

# 2. Problem Statement
Currently, accessing tenant-related information such as KYC/KYB status requires direct calls to third-party portals like AiRev each time a resource or page is accessed. This on-demand approach introduces latency and reliability concerns, especially if the external service is slow or temporarily unavailable. 

Moreover, tenant-related logic is fragmented across multiple areas in the codebase, making it difficult to maintain, scale, or extend for new clients without significant rework. These challenges collectively hinder our ability to deliver a consistent, real-time experience and complicate the onboarding of new client integrations.  

# 3. Proposed Solution
As part of our integration with third-party portals to retrieve tenant-related information such as KYC and KYB status. I propose the following solution: 

* **Dedicated repository** -To manage tenant-related data for all clients. While this will initially support Core42, the architecture will be designed to accommodate additional clients in the future, such as Airtel, with minimal changes. This repository will include modular "managers" that listen to changes in the Tenant and Tenant Users collections. Each manager will handle orchestration of external API calls based on tenant-related events. 

 

* **Polling Mechanism** -We cannot rely on a push-based or event-driven model from all clients (e.g., some might update KYC values periodically without explicit triggers like order creation). 

   This polling mechanism is necessary because:  

    We cannot rely on a push-based or event-driven model from all clients (e.g., some might update KYC values periodically without explicit triggers like order creation). 

    Some clients may provide passive updates, and our system needs to stay in sync by proactively fetching the latest information. 

```




                                                      Updating Coredge via API provided by 
                                                      Coredge whenever there is an update made 
                                                      to KYC/KYB or payment Configuration

                                       ┌◄────────────────────────────────────────────────────────────────┐ 
                                       │                                                                 │
                                       │                                                                 │
                           ┌───────────▼─────────┐                                              ┌────────▼───────────┐        
                           │                     │                                              │                    │        
                           │  Coredge's Platform ├◄───────────────────────────────────────────► │   AiRev's Platform │        
                           │                     │                                              │                    │         
                           └─────────────────────┘                                              └────────────────────┘        
                                                        Polling AiRev's GET tenant API to
                                                        make sure information is always updated 
                                                        in case AiRev  was not able to update us

```

## Configuration design
```bash
client:
  kyc:                         # Triggered on tenant user creation (individual admin case)
    enabled: bool
    server-path: string
    pollingTime: integer
    stopOnceSet: bool

  kyb:                         # Triggered on tenant creation
    enabled: bool
    server-path: string
    pollingTime: integer
    stopOnceSet: bool

  paymentMethodConfigured:     # Triggered on tenant creation
    enabled: bool
    server-path: string
    pollingTime: integer
    stopOnceSet: bool

  tenantType:                  # Triggered on tenant creation
    enabled: bool
    server-path: string
    pollingTime: integer
    stopOnceSet: bool  # Optional
```

This configuration Design enables easy extensibility — new data constructs can be added in the future simply by updating the config and toggling their enabled status. 

## Project Overview
```bash 

/shared
  └── manager/             Common manager logic, dynamically updated based on client name
/clients
  └── client_a/            Specific mappings and extensions for Core42
  └── client_b/            Specific mappings and extensions for Airtel
/cmd
  └── service/             Entry point for running the service

```  
This project structure separates shared logic from client-specific customizations for clarity and scalability: 

**/shared/manager/:** Contains the core manager logic used across all clients. It adapts dynamically based on the client's name, ensuring reusable and centralized functionality. 

**/clients/client_a/ & /client_b/:** Hold client-specific configurations and mappings — for example, Core42 and Airtel — allowing tailored behavior without modifying shared code. 

**/cmd/service/:** Serves as the main entry point for running the polling service, initializing managers and starting client-specific integrations. 

# 4.Alternatives considered 

# 5.API schema changes 

# 6. UI changes 

# 7. Notification impact 

# 8. Provisioning changes 

# 9. Implementation 

 * Coredge will be providing AiRev with specific APIs which will be used to update information at Coredge’s end and triggered by AiRev whenever there is a change in KYC status, KYB status or Payment Method configuration. 

 * AiRev would have to keep some simple retry mechanism also, if possible, in case one of the API hits fails initially. 

 * Coredge would not be completely dependent on AiRev to update information related to KYC/KYB or payment method configuration. 

 * We will be making use of polling as a backup mechanism to make sure that Coredge always has updated information at their end. 

 ```
                                                                           ┌──────────────┐                                                  
                                             ─◄────────────────────────────┼   AiRev      │                                                  
                                             │                             │              │                                                  
                                             │ Triggers API to update      └────────▲─────┘                                                  
                                             │ any status change in                 │                                                        
                                             │ KYC/KYB or PMC                       │                                                        
                                             │                              ┌───────▼─────┐                                                  
                                             │                              │ Manager     │ Polling GET Tenant API to                        
                                             │                              │             │  fetch KYC/KYB and PMC status                    
                                             │                              └──────▲──────┘  to make sure Coredge updated                    
                                             │                                     │                                                         
                                             │  ┌───────────┐                  ┌───▼──┐                                                      
                                             │  │           ┼◄────────────────►│      │                                                      
                                             │  │ Local     │Will store        │Mongo │◄───────────────────◄─────────────▲                   
                                             │  │ cache     │payment method    │ DB   │                   ▲              │  Fetching/Updating
                                             │  │           │configured or     │      │                   │              │  information in DB
                                             │  │           │not               └──────┘                   │              │                   
                                             │  │           │                                             │              │                   
                                             │  └───────────┘                                   ┌─────────┴─────┐        │                   
                                             │        ▲                ┌───────────────────────►┤Microservice 1 │        │                   
                                             │        │                │                        │               │        │                   
                                             ▼        ▼                │                        └───────────────┘        │                   
                                         ┌──────────────────┐          │                        ┌───────────────┐        │                   
┌────────────┐      ┌─────────────┐      │  Auth-Gateway    ├──────────┘                        │Microservice 2 │        │                   
│ Client     │      │ Frontend    │      │  Authentication, ┼──────────────────────────────────►┤               │ ──────►▲                   
│            ├─────►│             ├─────►│Authorization and ├──────────┐                        └───────────────┘        │                   
└────────────┘      └─────────────┘      │ redirect         │          │                        ┌───────────────┐        │                   
                                         └──────────────────┘          └───────────────────────►┤Microservice 3 │        │                   
                                         Auth-Gateway will also                                 │               ├───────►│                   
                                         Check for every request                                └───────────────┘                            
                                         whether payment method is                                                                           
                                         configured or not for                                                                               
                                          users's tenant                                                                                     
```

# 10. Performance and scaling impact 

**Scalability:** When a new client comes onboard with a similar requirement, we only need to update a shared plugin. No duplication across client-specific repositories. 

**Maintainability:** Shared logic is centralized, reducing redundancy and chances of inconsistency across implementations. 

**Flexibility:** New client-specific behaviors can be implemented easily via isolated configurations and client-specific plugins. 

 

# 11. Upgrade 

# 12. Deprecations 

# 13. Dependencies 

# 14. Security Considerations

# 15. Testing 

           Unit tests 

           Dev tests 

           System tests 

# 16. Documentation Impact 

# 17. References 

 

 