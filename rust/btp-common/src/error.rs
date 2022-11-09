pub mod errors {
    use serde::{Deserialize, Serialize};
    use std::fmt::{self, Error, Formatter};

    pub trait Exception {
        fn code(&self) -> u32;
        fn message(&self) -> String;
    }

    #[derive(Debug, PartialEq, Eq, Serialize, Deserialize)]
    pub enum BtpException<T: Exception> {
        Base,
        Bmc(T),
        Bmv(T),
        Bsh(T),
        Reserved,
    }

    impl<T> Exception for BtpException<T>
    where
        T: Exception,
    {
        fn code(&self) -> u32 {
            match self {
                BtpException::Base => 0,
                BtpException::Bmc(error) => error.code() + 10,
                BtpException::Bsh(error) => error.code() + 40,
                _ => todo!(),
            }
        }

        fn message(&self) -> String {
            match self {
                BtpException::Base => todo!(),
                BtpException::Bmc(error) => error.message(),
                BtpException::Bsh(error) => error.message(),
                _ => todo!(),
            }
        }
    }

    impl From<(u32, &Option<String>)> for Box<dyn Exception> {
        fn from((code, message): (u32, &Option<String>)) -> Self {
            match code {
                0..=9 => todo!(),
                10..=24 => Box::new(BtpException::Bmc(BmcError::from((code - 10, message)))),
                25..=39 => todo!(),
                _ => Box::new(BtpException::Bsh(BshError::from((code - 40, message)))),
            }
        }
    }

    #[derive(Debug, Clone, Eq, PartialEq, Serialize, Deserialize)]
    #[serde(tag = "error")]
    pub enum BmvError {
        Unknown { message: String },
        NotBmc,
        InvalidWitnessOld { message: String },
        InvalidWitnessNewer { message: String },
        InvalidBlockProof { message: String },
        InvalidVotes { message: String },
        InvalidBlockProofHeightHigher { expected: u64, actual: u64 },
        InvalidBlockUpdate { message: String },
        InvalidBlockUpdateHeightLower { expected: u64, actual: u64 },
        InvalidBlockUpdateHeightHigher { expected: u64, actual: u64 },
        DecodeFailed { message: String },
        EncodeFailed { message: String },
        InvalidReceipt { message: String },
        InvalidReceiptProof { message: String },
        InvalidEventProof { message: String },
        InvalidEventLog { message: String },
        InvalidSequence { expected: u128, actual: u128 },
        InvalidSequenceHigher { expected: u128, actual: u128 },
    }

    impl Exception for BmvError {
        fn code(&self) -> u32 {
            u32::from(self)
        }
        fn message(&self) -> String {
            self.to_string()
        }
    }

    impl From<&BmvError> for u32 {
        fn from(bsh_error: &BmvError) -> Self {
            match bsh_error {
                BmvError::Unknown { message: _ } => 0,
                _ => 0,
            }
        }
    }

    impl From<(u32, &Option<String>)> for BshError {
        fn from((code, _): (u32, &Option<String>)) -> BshError {
            match code {
                1 => BshError::PermissionNotExist,
                _ => BshError::Unknown,
            }
        }
    }

    impl From<(u32, &Option<String>)> for BmcError {
        fn from((code, message): (u32, &Option<String>)) -> BmcError {
            match code {
                0 => BmcError::Unknown {
                    message: message.clone().unwrap_or("Undefined".to_string()),
                },
                1 => BmcError::PermissionNotExist,
                2 => BmcError::InvalidSerialNo,
                3 => BmcError::VerifierExist,
                4 => BmcError::VerifierNotExist,
                5 => BmcError::ServiceExist,
                6 => BmcError::ServiceNotExist,
                7 => BmcError::LinkExist,
                8 => BmcError::LinkNotExist,
                9 => BmcError::RelayExist {
                    link: message.clone().unwrap_or("Undefined".to_string()),
                },
                10 => BmcError::RelayNotExist {
                    link: message.clone().unwrap_or("Undefined".to_string()),
                },
                11 => BmcError::Unreachable {
                    destination: message.clone().unwrap_or("Undefined".to_string()),
                },
                12 => BmcError::ErrorDrop,
                13 => BmcError::InvalidSequence,
                _ => BmcError::Unknown {
                    message: message.clone().unwrap_or("Undefined".to_string()),
                },
            }
        }
    }

    impl fmt::Display for BmvError {
        fn fmt(&self, f: &mut Formatter<'_>) -> Result<(), Error> {
            let label = "BMVRevert";
            match self {
                BmvError::NotBmc => {
                    write!(f, "{}{}", label, "NotBMC")
                }
                BmvError::DecodeFailed { message } => {
                    write!(f, "{}{}: {}", label, "DecodeError", message)
                }
                BmvError::EncodeFailed { message } => {
                    write!(f, "{}{}: {}", label, "EncodeError", message)
                }
                BmvError::InvalidBlockUpdate { message } => {
                    write!(f, "{}{}: {}", label, "InvalidBlockUpdate", message)
                }
                BmvError::InvalidVotes { message } => {
                    write!(f, "{}{}: {}", label, "InvalidVotes", message)
                }
                BmvError::InvalidBlockUpdateHeightLower { expected, actual } => {
                    write!(
                        f,
                        "{}{} expected: {}, but got: {}",
                        label, "InvalidBlockUpdateHeightLower", expected, actual
                    )
                }
                BmvError::InvalidBlockUpdateHeightHigher { expected, actual } => {
                    write!(
                        f,
                        "{}{} expected: {}, but got: {}",
                        label, "InvalidBlockUpdateHeightLower", expected, actual
                    )
                }
                BmvError::InvalidBlockProofHeightHigher { expected, actual } => {
                    write!(
                        f,
                        "{}{} expected: {}, but got: {}",
                        label, "InvalidBlockProofHeightHigher", expected, actual
                    )
                }
                BmvError::InvalidWitnessOld { message } => {
                    write!(f, "{}{}: {}", label, "InvalidWitnessOld", message)
                }
                BmvError::InvalidWitnessNewer { message } => {
                    write!(f, "{}{}: {}", label, "InvalidWitnessNewer", message)
                }
                BmvError::Unknown { message } => {
                    write!(f, "{}{}: {}", label, "Unknown", message)
                }
                BmvError::InvalidBlockProof { message } => {
                    write!(f, "{}{}: {}", label, "InvalidBlockProof", message)
                }
                BmvError::InvalidReceipt { message } => {
                    write!(f, "{}{}: {}", label, "InvalidReceipt", message)
                }
                BmvError::InvalidReceiptProof { message } => {
                    write!(f, "{}{}: {}", label, "InvalidReceiptProof", message)
                }
                BmvError::InvalidEventLog { message } => {
                    write!(f, "{}{}: {}", label, "InvalidEventLog", message)
                }
                BmvError::InvalidEventProof { message } => {
                    write!(f, "{}{}: {}", label, "InvalidEvenProof", message)
                }
                BmvError::InvalidSequence { expected, actual } => {
                    write!(
                        f,
                        "{}{} expected: {}, but got: {}",
                        label, "InvalidSequence", expected, actual
                    )
                }
                BmvError::InvalidSequenceHigher { expected, actual } => {
                    write!(
                        f,
                        "{}{} expected: {}, but got: {}",
                        label, "InvalidSequenceHigher", expected, actual
                    )
                }
            }
        }
    }

    #[derive(Debug, Clone, Eq, PartialEq, Serialize, Deserialize)]
    #[serde(tag = "error")]
    pub enum BshError {
        Unknown,
        LastOwner,
        OwnerExist,
        OwnerNotExist,
        PermissionNotExist,
        NotMinimumDeposit,
        NotMinimumRefundable,
        NotMinimumAmount,
        NotMinimumBalance { account: String },
        TokenExist,
        TokenNotExist { message: String },
        Failure,
        Reverted { message: String },
        NotBmc,
        InvalidService,
        DecodeFailed { message: String },
        EncodeFailed { message: String },
        InvalidSetting,
        InvalidCount { message: String },
        InvalidAddress { message: String },
        SameSenderReceiver,
        AccountNotExist,
        TokenNotRegistered,
        LessThanZero,
        UserAlreadyBlacklisted,
        BlacklistedUsers { message: String },
        NonBlacklistedUsers { message: String },
        InvalidParams,
        LimitExceed,
        LimitNotSet,
        RequiredMinimumOneYoctoNear,
    }

    impl Exception for BshError {
        fn code(&self) -> u32 {
            u32::from(self)
        }
        fn message(&self) -> String {
            self.to_string()
        }
    }

    impl From<&BshError> for u32 {
        fn from(bsh_error: &BshError) -> Self {
            match bsh_error {
                BshError::Unknown => 0,
                BshError::PermissionNotExist => 1,
                _ => 0,
            }
        }
    }

    impl fmt::Display for BshError {
        fn fmt(&self, f: &mut Formatter<'_>) -> Result<(), Error> {
            let label = "BSHRevert";
            match self {
                BshError::Reverted { message } => {
                    write!(f, "{}{}: {}", label, "Reverted", message)
                }
                BshError::TokenExist => {
                    write!(f, "{}{}", label, "AlreadyExistsToken")
                }
                BshError::TokenNotExist { message } => {
                    write!(f, "{}{}: {}", label, "NotExistsToken", message)
                }
                BshError::LastOwner => {
                    write!(f, "{}{}", label, "LastOwner")
                }
                BshError::OwnerExist => {
                    write!(f, "{}{}", label, "AlreadyExistsOwner")
                }
                BshError::OwnerNotExist => {
                    write!(f, "{}{}", label, "NotExistsOwner")
                }
                BshError::PermissionNotExist => {
                    write!(f, "{}{}", label, "NotExistsPermission")
                }
                BshError::NotMinimumDeposit => {
                    write!(f, "{}{}", label, "NotMinimumDeposit")
                }
                BshError::NotMinimumRefundable => {
                    write!(f, "{}{}", label, "NotMinimumRefundable")
                }
                BshError::NotBmc => {
                    write!(f, "{}{}", label, "NotBMC")
                }
                BshError::InvalidService => {
                    write!(f, "{}{}", label, "InvalidSvc")
                }
                BshError::DecodeFailed { message } => {
                    write!(f, "{}{} for {}", label, "DecodeError", message)
                }
                BshError::EncodeFailed { message } => {
                    write!(f, "{}{} for {}", label, "EncodeError", message)
                }
                BshError::InvalidSetting => {
                    write!(f, "{}{}", label, "InvalidSetting")
                }
                BshError::InvalidAddress { message } => {
                    write!(f, "{}{}: {}", label, "InvalidAddress", message)
                }
                BshError::InvalidCount { message } => {
                    write!(f, "{}{} for {}", label, "InvalidCount", message)
                }
                BshError::SameSenderReceiver => {
                    write!(f, "{}{}", label, "SameSenderReceiver")
                }
                BshError::AccountNotExist => {
                    write!(f, "{}{}", label, "AccountNotExist")
                }
                BshError::NotMinimumBalance { account } => {
                    write!(f, "{}{} for {}", label, "NotMinimumBalance", account)
                }
                BshError::NotMinimumAmount => {
                    write!(f, "{}{}", label, "NotMinimumAmount")
                }
                BshError::TokenNotRegistered => {
                    write!(f, "{}{}", label, "TokenNotRegistered")
                }
                BshError::Unknown => {
                    write!(f, "{}{}", label, "Unknown")
                }
                BshError::LessThanZero => {
                    write!(f, "{}{}", label, "LessThanZero")
                }
                BshError::Failure => {
                    write!(f, "{}{}", label, "Failure")
                }
                BshError::UserAlreadyBlacklisted => {
                    write!(f, "{}{}", label, "AlreadyBlacklisted")
                }
                BshError::NonBlacklistedUsers { message } => {
                    write!(f, "{}{} for {}", label, "UsersNotBlacklisted", message)
                }
                BshError::InvalidParams => {
                    write!(f, "{}{}", label, "InvalidParams")
                }
                BshError::LimitExceed => {
                    write!(f, "{}{}", label, "LimitExceed")
                }
                BshError::BlacklistedUsers { message } => {
                    write!(f, "{}{} for {}", label, "UsersBlacklisted", message)
                }
                BshError::LimitNotSet => {
                    write!(f, "{}{}", label, "LimitNotSet")
                }
                BshError::RequiredMinimumOneYoctoNear => {
                    write!(f, "{}{}", label, "RequiredMinimumOneYoctoNear")
                }
            }
        }
    }

    #[derive(Debug, Clone, Eq, PartialEq, Serialize, Deserialize)]
    #[serde(tag = "error")]
    pub enum BmcError {
        DecodeFailed { message: String },
        EncodeFailed { message: String },
        ErrorDrop,
        FeeAggregatorNotAllowed { source: String },
        InternalServiceCallNotAllowed { source: String },
        InvalidAddress { description: String },
        InvalidParam,
        InvalidSerialNo,
        LastOwner,
        LinkExist,
        LinkNotExist,
        LinkRouteExist,
        OwnerExist,
        OwnerNotExist,
        PermissionNotExist,
        RelayExist { link: String },
        RelayNotExist { link: String },
        RequestExist,
        RequestNotExist,
        RouteExist,
        RouteNotExist,
        ServiceExist,
        ServiceNotExist,
        Unknown { message: String },
        Unreachable { destination: String },
        VerifierExist,
        VerifierNotExist,
        Unauthorized { message: &'static str },
        InvalidSequence,
        InternalEventHandleNotExists,
        UnknownHandleBtpError,
        UnkownHandleBtpMessage,
    }

    impl Exception for BmcError {
        fn code(&self) -> u32 {
            u32::from(self)
        }
        fn message(&self) -> String {
            self.to_string()
        }
    }

    impl From<&BmcError> for u32 {
        fn from(bmc_error: &BmcError) -> Self {
            match bmc_error {
                BmcError::Unknown { message: _ } => 0,
                BmcError::PermissionNotExist => 1,
                BmcError::InvalidSerialNo => 2,
                BmcError::VerifierExist => 3,
                BmcError::VerifierNotExist => 4,
                BmcError::ServiceExist => 5,
                BmcError::ServiceNotExist => 6,
                BmcError::LinkExist => 7,
                BmcError::LinkNotExist => 8,
                BmcError::RelayExist { link: _ } => 9,
                BmcError::RelayNotExist { link: _ } => 10,
                BmcError::Unreachable { destination: _ } => 11,
                BmcError::ErrorDrop => 12,
                BmcError::InvalidSequence => 13,
                _ => 0,
            }
        }
    }

    impl fmt::Display for BmcError {
        fn fmt(&self, f: &mut Formatter<'_>) -> Result<(), Error> {
            let label = "BMCRevert";
            match self {
                BmcError::InvalidAddress { description } => {
                    write!(f, "{}{}: {}", label, "InvalidAddress", description)
                }
                BmcError::RequestExist => write!(f, "{}{}", label, "RequestPending"),
                BmcError::RequestNotExist => write!(f, "{}{}", label, "NotExistRequest"),
                BmcError::ServiceExist => write!(f, "{}{}", label, "AlreadyExistsBSH"),
                BmcError::ServiceNotExist => write!(f, "{}{}", label, "NotExistBSH"),
                BmcError::PermissionNotExist => write!(f, "{}{}", label, "NotExistsPermission"),
                BmcError::LastOwner => write!(f, "{}{}", label, "LastOwner"),
                BmcError::OwnerExist => write!(f, "{}{}", label, "AlreadyExistsOwner"),
                BmcError::OwnerNotExist => write!(f, "{}{}", label, "NotExistsOwner"),
                BmcError::LinkExist => write!(f, "{}{}", label, "AlreadyExistsLink"),
                BmcError::LinkNotExist => write!(f, "{}{}", label, "NotExistsLink"),
                BmcError::LinkRouteExist => write!(f, "{}{}", label, "LinkRouteExist"),
                BmcError::RouteExist => write!(f, "{}{}", label, "AlreadyExistsRoute"),
                BmcError::RouteNotExist => write!(f, "{}{}", label, "NotExistsRoute"),
                BmcError::InvalidParam => write!(f, "{}{}", label, "InvalidParam"),
                BmcError::VerifierExist => write!(f, "{}{}", label, "AlreadyExistsBMV"),
                BmcError::VerifierNotExist => write!(f, "{}{}", label, "NotExistBMV"),
                BmcError::InvalidSequence => write!(f, "{}{}", label, "InvalidSequence"),
                BmcError::RelayExist { link } => {
                    write!(f, "{}{} for {}", label, "RelayExist", link)
                }
                BmcError::RelayNotExist { link } => {
                    write!(f, "{}{} for {}", label, "NotExistRelay", link)
                }
                BmcError::DecodeFailed { message } => {
                    write!(f, "{}{} for {}", label, "DecodeError", message)
                }
                BmcError::EncodeFailed { message } => {
                    write!(f, "{}{} for {}", label, "EncodeError", message)
                }
                BmcError::ErrorDrop => {
                    write!(f, "{}{}", label, "ErrorDrop")
                }
                BmcError::InternalServiceCallNotAllowed { source } => {
                    write!(
                        f,
                        "{}{} for {}",
                        label, "NotAllowedInternalServiceCall", source
                    )
                }
                BmcError::FeeAggregatorNotAllowed { source } => {
                    write!(f, "{}{} from {}", label, "NotAllowedFeeAggregator", source)
                }
                BmcError::Unreachable { destination } => {
                    write!(f, "{}{} at {}", label, "Unreachable", destination)
                }
                BmcError::Unknown { message } => {
                    write!(f, "{}{}:{}", label, "Unknown", message)
                }
                BmcError::InvalidSerialNo => {
                    write!(f, "{}{}", label, "Invalid Serial No")
                }
                BmcError::Unauthorized { message } => {
                    write!(f, "{}{}: {}", label, "Unauthorized", message)
                }
                BmcError::InternalEventHandleNotExists => {
                    write!(f, "{}{}", label, "NotExistInternalEventHandle")
                }

                BmcError::UnknownHandleBtpError => {
                    write!(f, "{}{}", label, "UnknownHandleBtpError")
                }
                BmcError::UnkownHandleBtpMessage => {
                    write!(f, "{}{}", label, "UnkownHandleBtpMessage")
                }
            }
        }
    }
}
