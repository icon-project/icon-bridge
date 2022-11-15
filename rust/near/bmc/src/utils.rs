/// Trait utils
pub trait Utils {
    ///creates ceil div function
    fn ceil_div(num1: u128, num2: u128) -> u128;
    ///creates getScale function
    fn get_scale(block_interval_src: u128, block_interval_dst: u128) -> u128;
    ///creates getRotateTerm function
    fn get_rotate_term(max_agg: u128, scale: u128) -> u128;
}

impl Utils for u128 {

    ///
    ///ceil div function is created
    /// # Arguments
    /// * num1 and num2 should be given as arguments which are unsigned integers
    /// Returns the unsigned number
    /// 
    
    fn ceil_div(num1: u128, num2: u128) -> u128 {
        if num1 % num2 == 0 {
            return (num1 / num2) + 1;
        }
        (num1 / num2) + 1
    }


    ///
    /// get scale function is created
    /// # Arguments
    /// * Block_interval_src and block_interval_dst used as arguments which are unsigned integers
    /// Returns the unsigned number
    /// 

    fn get_scale(block_interval_src: u128, block_interval_dst: u128) -> u128 {
        Self::ceil_div(block_interval_src * 10_u128.pow(6), block_interval_dst)
    }

    ///
    /// get rotate term function is created
    /// #Arguments
    /// * `max_agg` - unsigned number
    /// * `scale`- unsigned number
    /// Returns the unsigned number
    /// 
    fn get_rotate_term(max_agg: u128, scale: u128) -> u128 {
        if scale > 0 {
            return Self::ceil_div(max_agg * 10_u128.pow(6), scale);
        }
        0
    }
}
