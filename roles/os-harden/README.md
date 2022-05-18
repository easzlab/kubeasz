# devsec.os_hardening

![devsec.os_hardening](https://github.com/dev-sec/ansible-os-hardening/workflows/devsec.os_hardening/badge.svg)

## Looking for the old ansible-os-hardening role?

This role is now part of the hardening-collection. You can find the old role in the branch `legacy`.

## Description

This role provides numerous security-related configurations, providing all-round base protection. It is intended to be compliant with the [DevSec Linux Baseline](https://github.com/dev-sec/linux-baseline).

It configures:

- Remove unused yum repositories and enable GPG key-checking
- Remove packages with known issues
- Configures pam for strong password checks
- Installs and configures auditd
- Disable core dumps via soft limits
- sets a restrictive umask
- Configures execute permissions of files in system paths
- Hardens access to shadow and passwd files
- Disables unused filesystems
- Disables rhosts
- Configures secure ttys
- Configures kernel parameters via sysctl
- Enables selinux on EL-based systems
- Remove SUIDs and GUIDs
- Configures login and passwords of system accounts

It will not:

- Update system packages
- Install security patches

## Requirements

- Ansible 2.9.0

## Known Limitations

### Docker support

If you're using Docker / Kubernetes+Docker you'll need to override the ipv4 ip forward sysctl setting.

```yaml
- hosts: localhost
  collections:
    - devsec.hardening
  roles:
    - devsec.hardening.os_hardening
  vars:
    sysctl_overwrite:
      # Enable IPv4 traffic forwarding.
      net.ipv4.ip_forward: 1
```

### sysctl - vm.mmap_rnd_bits

We are setting this sysctl to a default of `32`, some systems only support smaller values and this will generate an error. Unfortunately we cannot determine the correct applicable maximum. If you encounter this error you have to override this sysctl in your playbook.

```yaml
- hosts: localhost
  collections:
    - devsec.hardening
  roles:
    - devsec.hardening.os_hardening
  vars:
    sysctl_overwrite:
      vm.mmap_rnd_bits: 16
```

### Testing with inspec

If you're using inspec to test your machines after applying this role, please make sure to add the connecting user to the `os_ignore_users`-variable.
Otherwise inspec will fail. For more information, see [issue #124](https://github.com/dev-sec/ansible-os-hardening/issues/124).

We know that this is the case on Raspberry Pi.

## Variables

- `os_desktop_enable`
  - Default: `false`
  - Description: true if this is a desktop system, ie Xorg, KDE/GNOME/Unity/etc
- `os_env_extra_user_paths`
  - Default: `[]`
  - Description: add additional paths to the user's `PATH` variable (default is empty).
- `os_env_umask`
  - Default: `027`
  - Description: set default permissions for new files to `750`
- `os_auth_pw_max_age`
  - Default: `60`
  - Description: maximum password age (set to `99999` to effectively disable it)
- `os_auth_pw_min_age`
  - Default: `7`
  - Description: minimum password age (before allowing any other password change)
- `os_auth_retries`
  - Default: `5`
  - Description: the maximum number of authentication attempts, before the account is locked for some time
- `os_auth_lockout_time`
  - Default: `600`
  - Description: time in seconds that needs to pass, if the account was locked due to too many failed authentication attempts
- `os_auth_timeout`
  - Default: `60`
  - Description: authentication timeout in seconds, so login will exit if this time passes
- `os_auth_allow_homeless`
  - Default: `false`
  - Description: true if to allow users without home to login
- `os_auth_pam_passwdqc_enable`
  - Default: `true`
  - Description: true if you want to use strong password checking in PAM using passwdqc
- `os_auth_pam_passwdqc_options`
  - Default: `min=disabled,disabled,16,12,8`
  - Description: set to any option line (as a string) that you want to pass to passwdqc
- `os_security_users_allow`
  - Default: `[]`
  - Description: list of things, that a user is allowed to do. May contain `change_user`.
- `os_security_kernel_enable_module_loading`
  - Default: `true`
  - Description: true if you want to allowed to change kernel modules once the system is running (eg `modprobe`, `rmmod`)
- `os_security_kernel_enable_core_dump`
  - Default: `false`
  - Description: kernel is crashing or otherwise misbehaving and a kernel core dump is created
- `os_security_suid_sgid_enforce`
  - Default: `true`
  - Description: true if you want to reduce SUID/SGID bits. There is already a list of items which are searched for configured, but you can also add your own
- `os_security_suid_sgid_blacklist`
  - Default: `[]`
  - Description: a list of paths which should have their SUID/SGID bits removed
- `os_security_suid_sgid_whitelist`
  - Default: `[]`
  - Description: a list of paths which should not have their SUID/SGID bits altered
- `os_security_suid_sgid_remove_from_unknown`
  - Default: `false`
  - Description: true if you want to remove SUID/SGID bits from any file, that is not explicitly configured in a `blacklist`. This will make every Ansible-run search through the mounted filesystems looking for SUID/SGID bits that are not configured in the default and user blacklist. If it finds an SUID/SGID bit, it will be removed, unless this file is in your `whitelist`.
- `os_security_packages_clean`
  - Default: `true`
  - Description: removes packages with known issues. See section packages.
- `os_selinux_state`
  - Default: `enforcing`
  - Description: Set the SELinux state, can be either disabled, permissive, or enforcing.
- `os_selinux_policy`
  - Default: `targeted`
  - Description: Set the SELinux polixy.
- `ufw_manage_defaults`
  - Default: `true`
  - Description: true means apply all settings with `ufw_` prefix
- `ufw_ipt_sysctl`
  - Default: `''`
  - Description: by default it disables IPT_SYSCTL in /etc/default/ufw. If you want to overwrite /etc/sysctl.conf values using ufw - set it to your sysctl dictionary, for example `/etc/ufw/sysctl.conf`
- `ufw_default_input_policy`
  - Default: `DROP`
  - Description: set default input policy of ufw to `DROP`
- `ufw_default_output_policy`
  - Default: `ACCEPT`
  - Description: set default output policy of ufw to `ACCEPT`
- `ufw_default_forward_policy`
  - Default: `DROP`
  - Description: set default forward policy of ufw to `DROP`
- `os_auditd_enabled`
  - Default: `true`
  - Description: Set to false to disable installing and configuring auditd.
- `os_auditd_max_log_file_action`
  - Default: `keep_logs`
  - Description: Defines the behaviour of auditd when its log file is filled up. Possible other values are described in the auditd.conf man page. The most common alternative to the default may be `rotate`.
- `hidepid_option`
  - Default: `2`
  - Description: `0`: This is the default setting and gives you the default behaviour. `1`: With this option an normal user would not see other processes but their own about ps, top etc, but he is still able to see process IDs in /proc. `2`: Users are only able too see their own processes (like with hidepid=1), but also the other process IDs are hidden for them in /proc.
- `proc_mnt_options`
  - Default: `rw,nosuid,nodev,noexec,relatime,hidepid={{ hidepid_option }}`
  - Description: Mount proc with hardenized options, including `hidepid` with variable value.

## Packages

We remove the following packages:

- xinetd ([NSA](https://apps.nsa.gov/iaarchive/library/ia-guidance/security-configuration/operating-systems/guide-to-the-secure-configuration-of-red-hat-enterprise.cfm), Chapter 3.2.1)
- inetd ([NSA](https://apps.nsa.gov/iaarchive/library/ia-guidance/security-configuration/operating-systems/guide-to-the-secure-configuration-of-red-hat-enterprise.cfm), Chapter 3.2.1)
- tftp-server ([NSA](https://apps.nsa.gov/iaarchive/library/ia-guidance/security-configuration/operating-systems/guide-to-the-secure-configuration-of-red-hat-enterprise.cfm), Chapter 3.2.5)
- ypserv ([NSA](https://apps.nsa.gov/iaarchive/library/ia-guidance/security-configuration/operating-systems/guide-to-the-secure-configuration-of-red-hat-enterprise.cfm), Chapter 3.2.4)
- telnet-server ([NSA](https://apps.nsa.gov/iaarchive/library/ia-guidance/security-configuration/operating-systems/guide-to-the-secure-configuration-of-red-hat-enterprise.cfm), Chapter 3.2.2)
- rsh-server ([NSA](https://apps.nsa.gov/iaarchive/library/ia-guidance/security-configuration/operating-systems/guide-to-the-secure-configuration-of-red-hat-enterprise.cfm), Chapter 3.2.3)
- prelink ([open-scap](https://static.open-scap.org/ssg-guides/ssg-sl7-guide-ospp-rhel7-server.html#xccdf_org.ssgproject.content_rule_disable_prelink))

## Disabled filesystems

We disable the following filesystems, because they're most likely not used:

- "cramfs"
- "freevxfs"
- "jffs2"
- "hfs"
- "hfsplus"
- "squashfs"
- "udf"
- "vfat" # only if uefi is not in use

To prevent some of the filesystems from being disabled, add them to the `os_filesystem_whitelist` variable.

## Example Playbook

```yaml
- hosts: localhost
  collections:
    - devsec.hardening
  roles:
    - devsec.hardening.os_hardening
```

## Changing sysctl variables

If you want to override sysctl-variables, you can use the `sysctl_overwrite` variable (in older versions you had to override the whole `sysctl_dict`).
So for example if you want to change the IPv4 traffic forwarding variable to `1`, do it like this:

```yaml
- hosts: localhost
  collections:
    - devsec.hardening
  roles:
    - devsec.hardening.os_hardening
  vars:
    sysctl_overwrite:
      # Enable IPv4 traffic forwarding.
      net.ipv4.ip_forward: 1
```

Alternatively you can change Ansible's [hash-behaviour](https://docs.ansible.com/ansible/latest/reference_appendices/config.html#default-hash-behaviour) to `merge`, then you only have to overwrite the single hash you need to. But please be aware that changing the hash-behaviour changes it for all your playbooks and is not recommended by Ansible.

## Improving Kernel Audit logging

By default, any process that starts before the `auditd` daemon will have an AUID of `4294967295`. To improve this and provide more accurate logging, it's recommended to add the kernel boot parameter `audit=1` to you configuration. Without doing this, you will find that your `auditd` logs fail to properly audit all processes.

For more information, please see this [upstream documentation](https://www.kernel.org/doc/html/latest/admin-guide/kernel-parameters.html) and your system's boot loader documentation for how to configure additional kernel parameters.

## More information

This role is mostly based on guides by:

- [Arch Linux wiki, Sysctl hardening](https://wiki.archlinux.org/index.php/Sysctl)
- [NSA: Guide to the Secure Configuration of Red Hat Enterprise Linux 5](http://www.nsa.gov/ia/_files/os/redhat/rhel5-guide-i731.pdf)
- [Ubuntu Security/Features](https://wiki.ubuntu.com/Security/Features)
- [Deutsche Telekom, Group IT Security, Security Requirements (German)](https://www.telekom.com/psa)
